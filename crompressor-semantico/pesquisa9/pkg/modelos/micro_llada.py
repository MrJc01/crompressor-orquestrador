import torch
import torch.nn as nn
from torch.nn import functional as F
import math

# Parâmetros Base do Micro-LLaDA
class LLaDAConfig:
    def __init__(self):
        self.vocab_size = 32000
        self.n_embd = 256
        self.n_head = 4
        self.n_layer = 6
        self.block_size = 256
        self.dropout = 0.1
        self.mask_token_id = 1 # Do nosso BPE Tokenizer
        self.pad_token_id = 0

class BidirectionalSelfAttention(nn.Module):
    def __init__(self, config):
        super().__init__()
        assert config.n_embd % config.n_head == 0
        self.c_attn = nn.Linear(config.n_embd, 3 * config.n_embd)
        self.c_proj = nn.Linear(config.n_embd, config.n_embd)
        self.n_head = config.n_head
        self.n_embd = config.n_embd
        # Sem máscara causal (tril)! Isto é Full Attention bidirecional.
        
    def forward(self, x):
        B, T, C = x.size()
        qkv = self.c_attn(x)
        q, k, v = qkv.split(self.n_embd, dim=2)
        k = k.view(B, T, self.n_head, C // self.n_head).transpose(1, 2)
        q = q.view(B, T, self.n_head, C // self.n_head).transpose(1, 2)
        v = v.view(B, T, self.n_head, C // self.n_head).transpose(1, 2)

        # Full Attention
        att = (q @ k.transpose(-2, -1)) * (1.0 / math.sqrt(k.size(-1)))
        att = F.softmax(att, dim=-1)
        
        y = att @ v
        y = y.transpose(1, 2).contiguous().view(B, T, C)
        return self.c_proj(y)

class MLP(nn.Module):
    def __init__(self, config):
        super().__init__()
        self.c_fc    = nn.Linear(config.n_embd, 4 * config.n_embd)
        self.c_proj  = nn.Linear(4 * config.n_embd, config.n_embd)

    def forward(self, x):
        x = self.c_fc(x)
        x = F.gelu(x)
        x = self.c_proj(x)
        return x

class Block(nn.Module):
    def __init__(self, config):
        super().__init__()
        self.ln_1 = nn.LayerNorm(config.n_embd)
        self.attn = BidirectionalSelfAttention(config)
        self.ln_2 = nn.LayerNorm(config.n_embd)
        self.mlp = MLP(config)

    def forward(self, x):
        x = x + self.attn(self.ln_1(x))
        x = x + self.mlp(self.ln_2(x))
        return x

class MicroLLaDA(nn.Module):
    def __init__(self, config):
        super().__init__()
        self.config = config
        self.transformer = nn.ModuleDict(dict(
            wte = nn.Embedding(config.vocab_size, config.n_embd),
            wpe = nn.Embedding(config.block_size, config.n_embd),
            h = nn.ModuleList([Block(config) for _ in range(config.n_layer)]),
            ln_f = nn.LayerNorm(config.n_embd),
        ))
        self.lm_head = nn.Linear(config.n_embd, config.vocab_size, bias=False)
        self.transformer.wte.weight = self.lm_head.weight # Weight Tying

    def forward(self, idx):
        device = idx.device
        B, T = idx.size()
        pos = torch.arange(0, T, dtype=torch.long, device=device)
        
        tok_emb = self.transformer.wte(idx)
        pos_emb = self.transformer.wpe(pos)
        x = tok_emb + pos_emb
        
        for block in self.transformer.h:
            x = block(x)
            
        x = self.transformer.ln_f(x)
        logits = self.lm_head(x)
        return logits

    def apply_forward_masking(self, batch_tokens):
        """
        Substitui aleatoriamente t% dos tokens pelo ID do [MASK].
        """
        B, T = batch_tokens.size()
        device = batch_tokens.device
        
        # Sorteia t ~ U[0, 1] para cada batch
        t = torch.rand(B, 1, device=device)
        
        # Probabilidade de ser mascarado (mantemos 0.0 probabilidade para paddings se necessário)
        prob_matrix = torch.full((B, T), 0.0, device=device) + t
        mask = torch.bernoulli(prob_matrix).bool()
        
        masked_inputs = batch_tokens.clone()
        masked_inputs[mask] = self.config.mask_token_id
        return masked_inputs, mask

    @torch.no_grad()
    def reverse_process(self, seq_len, num_steps=20, device='cpu'):
        """
        Denoising Iterativo (Inferência)
        """
        # Começa com matriz 100% de ruído (MASK)
        seq = torch.full((1, seq_len), self.config.mask_token_id, device=device)
        
        for step in range(num_steps):
            logits = self.forward(seq)
            probs = F.softmax(logits, dim=-1)
            
            # Pega o token mais provável e a sua confiança
            confianca, tokens_preditos = torch.max(probs, dim=-1)
            
            # O rácio de tokens para REVELAR cresce linearmente
            ratio_revelar = (step + 1) / num_steps
            k = int(seq_len * ratio_revelar)
            
            if k == seq_len:
                seq = tokens_preditos
                break
                
            # Manter os 'k' tokens mais confiantes
            limite_confianca = torch.topk(confianca, k).values[0, -1]
            mascara_rejeicao = confianca < limite_confianca
            
            seq = tokens_preditos.clone()
            seq[mascara_rejeicao] = self.config.mask_token_id
            
        return seq

if __name__ == "__main__":
    print("[*] Instanciando Micro-LLaDA CROM (Pesquisa 9)...")
    config = LLaDAConfig()
    model = MicroLLaDA(config)
    total_params = sum(p.numel() for p in model.parameters())
    print(f"[+] Modelo criado com {total_params / 1e6:.2f}M parâmetros.")
    
    # Teste de Treino Dummy
    texto_original = torch.randint(2, 100, (4, 32)) # (Batch, SeqLen)
    masked_inputs, mask = model.apply_forward_masking(texto_original)
    
    print(f"[*] Input Original Exemplo: {texto_original[0][:10].tolist()}")
    print(f"[*] Mascarado Exemplo (ID 1): {masked_inputs[0][:10].tolist()}")
    
    optimizer = torch.optim.AdamW(model.parameters(), lr=1e-3)
    
    # Loss convergence check
    for epoch in range(10):
        logits = model(masked_inputs)
        # Loss é calculada apenas onde mascarámos
        # Flatten para B*T
        logits_flat = logits.view(-1, config.vocab_size)
        targets_flat = texto_original.view(-1)
        mask_flat = mask.view(-1)
        
        loss = F.cross_entropy(logits_flat[mask_flat], targets_flat[mask_flat])
        
        optimizer.zero_grad()
        loss.backward()
        optimizer.step()
        
        print(f"Época {epoch} | Loss (só nas máscaras): {loss.item():.4f}")
        
    print("[*] Teste de Reverse Process (Inferência)...")
    out = model.reverse_process(seq_len=10, num_steps=5)
    print("Tokens gerados:", out[0].tolist())
