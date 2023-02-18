package handler

type Handler interface {
	validate() error
	core()
	Execute()
}

type BaseHandler struct {
	Handler Handler
}

func (h *BaseHandler) Execute() {
	err := h.Handler.validate()
	if err != nil {
		print("[Error] Validate")
		return
	}

	h.Handler.core()
}

type LoginHandler struct {
	BaseHandler
}

func (h *LoginHandler) validate() error {
	print("login validate\n")
	return nil
}

func (h *LoginHandler) core() {
	print("login core code\n")
}
