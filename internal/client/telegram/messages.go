package telegram

const (
	msgHelp = `
/start -- greet and display help!
/help -- display help
/add [Name] -- add task with name [Name] if not exists
/remove [Name] -- delete task with name [Name] if exists
/list -- list tasks
`
	msgHello           = `Hi there! ` + msgHelp
	msgUnknownCommand  = `Unknown command!`
	msgCreated         = `Created!`
	msgDeleted         = `Deleted!`
	msgAlreadyExists   = `Already exists!`
	msgNoTaskFound     = `No task found!`
	msgTaskDoesntExist = "Task doesnt exist!"
)
