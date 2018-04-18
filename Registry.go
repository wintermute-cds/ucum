package ucum

import "github.com/bertverhees/ucum/decimal"

type Registry struct {
	handlers map[string]SpecialUnitHandlerer
}

func NewRegistry() *Registry {
	r := &Registry{}
	r.handlers = make(map[string]SpecialUnitHandlerer)
	r.register(&CelsiusHandler{})
	r.register(&FahrenheitHandler{})
	r.register(NewHoldingHandler("[p'diop]", "deg", decimal.Zero))
	r.register(NewHoldingHandler("%[slope]", "deg", decimal.Zero))
	r.register(NewHoldingHandler("[hp_X]", "1", decimal.Zero))
	r.register(NewHoldingHandler("[hp_C]", "1", decimal.Zero))
	r.register(NewHoldingHandler("[pH]", "mol/l", decimal.Zero))
	r.register(NewHoldingHandler("Np", "1", decimal.Zero))
	r.register(NewHoldingHandler("B", "1", decimal.Zero))
	d, _ := decimal.NewFromString("2")
	r.register(NewHoldingHandler("B[SPL]", "10*-5.Pa", d))
	r.register(NewHoldingHandler("B[V]", "V", decimal.Zero))
	r.register(NewHoldingHandler("B[mV]", "mV", decimal.Zero))
	r.register(NewHoldingHandler("B[uV]", "uV", decimal.Zero))
	r.register(NewHoldingHandler("B[W]", "W", decimal.Zero))
	r.register(NewHoldingHandler("B[kW]", "kW", decimal.Zero))
	r.register(NewHoldingHandler("bit_s", "1", decimal.Zero))
	return r
}

func (r *Registry) register(handler SpecialUnitHandlerer) {
	r.handlers[handler.GetCode()] = handler
}

func (r *Registry) Exists(code string) bool {
	return r.handlers[code] != nil
}

func (r *Registry) Get(code string) SpecialUnitHandlerer {
	return r.handlers[code]
}
