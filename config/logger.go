package config

type Mail struct {
	Subject      string
	Sender       string
	Receivers    string
	MjAPIPublic  string
	MjAPIPrivate string
}

type Option struct {
	Directory string
	Filename  string
	Maxday    int64
	Priority  string
	Tag       string
}

type Hooks struct {
	Name    string
	Options Option
	Mail    Mail
}

type Out struct {
	Name    string
	Options Option
}

type Formatter struct {
	Name string
}

type LoggerConfiguration struct {
	Out       Out
	Level     string
	Formatter Formatter
	Hooks     Hooks
}
