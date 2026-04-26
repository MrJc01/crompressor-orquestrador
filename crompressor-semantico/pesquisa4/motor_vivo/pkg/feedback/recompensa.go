package feedback

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
)

type BrainState struct {
	Limiares map[string]float64 `json:"limiares"`
}

type MotorRecompensa struct {
	mu       sync.Mutex
	State    BrainState
	filepath string
}

func NovoMotorRecompensa(filepath string) *MotorRecompensa {
	mr := &MotorRecompensa{
		State: BrainState{
			Limiares: make(map[string]float64),
		},
		filepath: filepath,
	}
	mr.CarregarEstado()
	return mr
}

func (mr *MotorRecompensa) CarregarEstado() {
	if _, err := os.Stat(mr.filepath); os.IsNotExist(err) {
		// Inicializa com valor padrão caso não exista
		mr.State.Limiares["default"] = 4.0 // Ponto médio padrão da Sigmóide
		mr.SalvarEstado()
		return
	}

	dados, err := ioutil.ReadFile(mr.filepath)
	if err == nil {
		json.Unmarshal(dados, &mr.State)
	}
}

func (mr *MotorRecompensa) SalvarEstado() {
	mr.mu.Lock()
	defer mr.mu.Unlock()
	dados, _ := json.MarshalIndent(mr.State, "", "  ")
	ioutil.WriteFile(mr.filepath, dados, 0644)
}

// PunirBucket ajusta a rigidez do motor para um contexto específico.
// Se o utilizador disser "/errado", o motor reduz o Ponto Médio da sigmoide,
// tornando a deduplicação mais difícil e exigente (diminuindo falsos positivos).
func (mr *MotorRecompensa) PunirBucket(bucketID string) {
	mr.mu.Lock()
	val, existe := mr.State.Limiares[bucketID]
	if !existe {
		val = 4.0
	}
	// Diminui o ponto médio em 15% (tornando a sigmoide mais estrita)
	novoVal := val * 0.85
	if novoVal < 1.0 { // Limite mínimo de tolerância
		novoVal = 1.0
	}
	mr.State.Limiares[bucketID] = novoVal
	mr.mu.Unlock()
	mr.SalvarEstado()
}

func (mr *MotorRecompensa) ObterPontoMedio(bucketID string) float64 {
	mr.mu.Lock()
	defer mr.mu.Unlock()
	if val, existe := mr.State.Limiares[bucketID]; existe {
		return val
	}
	return 4.0 // Padrão
}
