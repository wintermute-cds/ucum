package ucum

type Registry struct {
	handlers map[string]SpecialUnitHandlerer
}

func NewRegistry() *Registry {
	r := &Registry{}
	r.handlers = make(map[string]SpecialUnitHandlerer)
	r.register(&CelsiusHandler{})
	r.register(&FahrenheitHandler{})
	r.register(NewHoldingHandler("[p'diop]", "deg", Zero))
	r.register(NewHoldingHandler("%[slope]", "deg", Zero))
	r.register(NewHoldingHandler("[hp_X]", "1", Zero))
	r.register(NewHoldingHandler("[hp_C]", "1", Zero))
	r.register(NewHoldingHandler("[pH]", "mol/l", Zero))
	r.register(NewHoldingHandler("Np", "1", Zero))
	r.register(NewHoldingHandler("B", "1", Zero))
	d, _ := NewDecimal("2")
	r.register(NewHoldingHandler("B[SPL]", "10*-5.Pa", d))
	r.register(NewHoldingHandler("B[V]", "V", Zero))
	r.register(NewHoldingHandler("B[mV]", "mV", Zero))
	r.register(NewHoldingHandler("B[uV]", "uV", Zero))
	r.register(NewHoldingHandler("B[W]", "W", Zero))
	r.register(NewHoldingHandler("B[kW]", "kW", Zero))
	r.register(NewHoldingHandler("bit_s", "1", Zero))
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
