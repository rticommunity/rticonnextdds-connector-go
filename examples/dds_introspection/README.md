# DDS Introspection
``` bash
./dds_introspection -h
NAME:
   DDS Introspection - Injecting DDS data for diagnostics

USAGE:
   dds_introspection [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --cfgPath value, -c value      Path for XML App Creation configuration (default: "ShapeExample.xml")
   --participant value, -p value  Participant name in XML (default: "MyParticipantLibrary::Zero")
   --writer value, -w value       Writer name in XML (default: "MyPublisher::MySquareWriter")
   --data value, -d value         Data in JSON format (default: "{\"color\": \"BLUE\", \"x\": 10, \"y\": 20, \"shapesize\": 30}")
   --count value                  Number of injecting samples (default: 10)
   --interval value               Interval between samples in seconds (default: 1)
   --help, -h                     show help
   --version, -v                  print the version
```
