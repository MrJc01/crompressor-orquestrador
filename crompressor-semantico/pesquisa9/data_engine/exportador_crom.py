import os
import struct
import torch
import numpy as np
import sys

# Adiciona o pacote de modelos ao path para poder importar a classe
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '..', 'pkg', 'modelos')))
from micro_llada import MicroLLaDA, LLaDAConfig

def exportar_modelo_crom(modelo, filepath):
    """
    Serializa os tensores de PyTorch num binário flat .crom para Go
    Formato para cada Tensor:
      [Int32] Tamanho do nome do tensor (L)
      [Bytes] Nome do tensor (tamanho L, utf-8)
      [Int32] N = número de dimensões
      [Int32] * N = shape do tensor
      [Float32] * (tamanho total do tensor) = dados puros do tensor
    """
    print(f"[*] A exportar a matriz cerebral para {filepath}...")
    state_dict = modelo.state_dict()
    
    # Criar um tensor flat
    with open(filepath, 'wb') as f:
        # Cabeçalho Mágico CROM
        f.write(b"CROM")
        # Número de tensores
        f.write(struct.pack('<I', len(state_dict)))
        
        for name, tensor in state_dict.items():
            # Nome
            name_bytes = name.encode('utf-8')
            f.write(struct.pack('<I', len(name_bytes)))
            f.write(name_bytes)
            
            # Shape
            shape = list(tensor.shape)
            f.write(struct.pack('<I', len(shape)))
            for s in shape:
                f.write(struct.pack('<I', s))
                
            # Dados Puros (Float32)
            tensor_np = tensor.detach().cpu().numpy().astype(np.float32)
            f.write(tensor_np.tobytes())
            
            print(f"    -> Exportado {name} | Shape: {shape}")
            
    print("[+] Cérebro exportado com sucesso. Pronto para a Borda (Go).")

if __name__ == "__main__":
    config = LLaDAConfig()
    modelo = MicroLLaDA(config)
    
    # Na Pesquisa 9 atual, o modelo inicializa com pesos aleatórios.
    # No pipeline de produção, rodaríamos modelo.load_state_dict() primeiro.
    # Aqui exportamos o estado atual para validar a matemática Go-Native.
    
    caminho_saida = "llada_10m.crom"
    exportar_modelo_crom(modelo, caminho_saida)
