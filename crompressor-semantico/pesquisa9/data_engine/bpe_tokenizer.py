import json
from collections import Counter

class BPETokenizer:
    def __init__(self, vocab_size=2000):
        self.vocab_size = vocab_size
        self.merges = {}
        self.vocab = {}
        # Tokens Especiais para Difusão
        self.SPECIAL_TOKENS = {"[PAD]": 0, "[MASK]": 1, "[UNK]": 2}
        
    def get_stats(self, ids):
        counts = Counter()
        for pair in zip(ids, ids[1:]):
            counts[pair] += 1
        return counts

    def merge(self, ids, pair, idx):
        newids = []
        i = 0
        while i < len(ids):
            if i < len(ids) - 1 and ids[i] == pair[0] and ids[i+1] == pair[1]:
                newids.append(idx)
                i += 2
            else:
                newids.append(ids[i])
                i += 1
        return newids

    def train(self, text):
        print(f"[*] Treinando BPE Tokenizer para vocab size: {self.vocab_size}...")
        tokens = list(text.encode("utf-8"))
        num_merges = self.vocab_size - 256 - len(self.SPECIAL_TOKENS)
        
        # Base UTF-8 bytes
        self.vocab = {idx: bytes([idx]) for idx in range(256)}
        
        for i, (special, idx) in enumerate(self.SPECIAL_TOKENS.items()):
            self.vocab[idx] = special.encode("utf-8")
        
        current_idx = 256 + len(self.SPECIAL_TOKENS)
        
        for i in range(num_merges):
            stats = self.get_stats(tokens)
            if not stats:
                break
            best_pair = max(stats, key=stats.get)
            tokens = self.merge(tokens, best_pair, current_idx)
            self.merges[best_pair] = current_idx
            # Simular a construção da string (apenas bytes concatenados)
            self.vocab[current_idx] = self.vocab[best_pair[0]] + self.vocab[best_pair[1]]
            current_idx += 1
            
        print("[+] Tokenizer treinado com sucesso!")
        
    def encode(self, text):
        tokens = list(text.encode("utf-8"))
        while len(tokens) >= 2:
            stats = self.get_stats(tokens)
            pair = min(stats, key=lambda p: self.merges.get(p, float("inf")))
            if pair not in self.merges:
                break
            idx = self.merges[pair]
            tokens = self.merge(tokens, pair, idx)
        return tokens

    def decode(self, ids):
        tokens = b"".join([self.vocab.get(idx, b"") for idx in ids])
        return tokens.decode("utf-8", errors="replace")

if __name__ == "__main__":
    corpus = "O universo é vasto e misterioso. A difusão transforma texto num vetor quântico. O CROM é rápido." * 10
    tok = BPETokenizer(vocab_size=300)
    tok.train(corpus)
    t = tok.encode("O universo é vasto")
    print("Tokens:", t)
    print("Decoded:", tok.decode(t))
