# gop
Process manager written in golang

Examples:

gop start <cmd> <args> //Starts a process in the manager (random instance name)
gop start <cmd> <args> -n <name> //Starts a named process in the manager

gop stop <instancename> //Starts a named process in the manager


gop start -s <cmd> <args> -n <name> //Starts a service that will launch on gop dameon launch. 
// The process instance name will be called the same as the service name

gop start -s -n <name> //Re-launch a service that isnt running

gop <enable/disable> <service> //Sets the enabled status of the service
gop remove <service> //removes a service


gop list //Get a list of all services and running process instances
gop log  //Get the logs for the dameons current session

Extra flags:
-v //voberse output
-r //Relitive path. Defaultly gop will change your cmd to its absolue path. This will stop that. (for things in the bin folder or for dameons not on the local machine)
