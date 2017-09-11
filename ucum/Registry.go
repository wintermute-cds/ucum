package ucum

type Registry struct{
	handlers map[string]SpecialUnitHandlerer
}

func NewRegistry()*Registry{
	r := &Registry{}
	r.handlers = make(map[string]SpecialUnitHandlerer)
	r.register(&CelsiusHandler{})
	r.register(&FahrenheitHandler{})
	r.register(NewHoldingHandler("[p'diop]", "deg",nil))
	r.register(NewHoldingHandler("%[slope]", "deg",nil))
	r.register(NewHoldingHandler("[hp_X]", "1",nil))
	r.register(NewHoldingHandler("[hp_C]", "1",nil))
	r.register(NewHoldingHandler("[pH]", "mol/l",nil))
	r.register(NewHoldingHandler("Np", "1",nil))
	r.register(NewHoldingHandler("B", "1",nil))
	d, _ := NewDecimal("2")
	r.register(NewHoldingHandler("B[SPL]", "10*-5.Pa", d))
	r.register(NewHoldingHandler("B[V]", "V",nil))
	r.register(NewHoldingHandler("B[mV]", "mV",nil))
	r.register(NewHoldingHandler("B[uV]", "uV",nil))
	r.register(NewHoldingHandler("B[W]", "W",nil))
	r.register(NewHoldingHandler("B[kW]", "kW",nil))
	r.register(NewHoldingHandler("bit_s", "1",nil))
	return r
}

func (r *Registry)register(handler SpecialUnitHandlerer){
	r.handlers[handler.GetCode()] = handler
}

func (r *Registry)Exists(code string)bool{
	return r.handlers[code] != nil
}

func (r *Registry)Get(code string)SpecialUnitHandlerer{
	return r.handlers[code]
}
