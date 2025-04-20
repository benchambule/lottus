package lottus

import (
	"errors"
)

func (app App) ProcRequest(req Request) (error, Message) {
	if app.initial == "" {
		return errors.New("Initial is not set"), Message{}
	}

	_, session := app.sessionStorage.Get(req.Msisdn)

	if session.GetLocation() == "" {
		_, msg := Redirect(app, app.initial, &req)

		session.Msisdn = req.Msisdn
		session.CurrentMessage = msg

		app.sessionStorage.Add(session)

		return nil, session.CurrentMessage
	}

	proc := app.Processor(session.GetLocation())

	// TODO: Verify whether location has displayer or processor

	message := proc(&req, session.CurrentMessage)
	session.CurrentMessage = message

	app.sessionStorage.Update(session)

	return nil, message
}

type InputValidationResult struct {
	UserInput   string
	NextMessage string
	Parameters  []Parameter
}

func Redirect(app App, loc string, req *Request) (error, Message) {
	proc := app.Displayer(loc)

	msg := proc(req, Message{})
	msg.Location = loc

	return nil, msg
}

type App struct {
	initial        string
	sessionStorage SessionStorage
	processors     map[string]MessageProcessor
}

func New(inital string, storage SessionStorage) App {
	return App{
		initial:        inital,
		sessionStorage: storage,
		processors:     make(map[string]MessageProcessor),
	}
}

type Processor func(r *Request, res Message) Message

func (app App) At(n string, d Processor, p Processor) {
	app.processors[n] = MessageProcessor{Displayer: d, Process: p}
}

func (app App) What(n string) (Displayer Processor, Process Processor) {
	mp := app.processors[n]

	return mp.Displayer, mp.Process
}

func (app App) Displayer(n string) Processor {
	p := app.processors[n]

	return p.Displayer
}

func (app App) Processor(n string) Processor {
	p := app.processors[n]

	return p.Process
}

type Session struct {
	Msisdn         string
	CurrentMessage Message
}

func (s Session) GetLocation() string {
	return s.CurrentMessage.Location
}

type SessionStorage interface {
	Get(i string) (error, Session)
	Add(s Session) error
	Delete(s Session) error
	Update(s Session) error
}

type Request struct {
	Msisdn     string
	Prompt     string
	Parameters []Parameter
}

type Parameter struct {
	Name  string
	Value any
}

type Input interface {
	Validate(r Request, msg Message) (error, InputValidationResult)
}

type Option struct {
	Key            string
	Text           string
	NextMessage    *Message
	ParameterValue string
}

type Message struct {
	Text       string
	Input      Input
	Parameters []Parameter
	Errors     []string
	Location   string
}
