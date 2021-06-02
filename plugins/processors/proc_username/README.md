# Procusername Processor Plugin

The Procusername processor plugin adds "username" tag to measurements which collected by win_perf_counters input plugin

It's useful to whom run multiple websites on one server, because the performance counter shows alike name such as "w3wp#1" for different instances which are started by same application, therefore we can't distinguish them by name. 

As we know, different app pools run using different  windows user. And generally speaking, we ought to run website using seprated app pools,  which means that, we could recognize them by processes' username.

## Effect Module

* inputs.win_perf_counters

## Effect Measurements

* win_proc
* win_dotnet

## Requirements

### Change perfmon process name format

|||
|  ----  | ----  |
| Key name | HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Services\PerfProc\Performance |
| Value name | ProcessNameFormat |
| Value type | REG_DWORD |
| Value | 2 (0x00000002) |

### Change dotnet perfmon process name format

|||
|  ----  | ----  |
| Key name | HKEY_LOCAL_MACHINE\System\CurrentControlSet\Services\.NETFramework\Performance|
| Value name | ProcessNameFormat |
| Value type | REG_DWORD |
| Value | 1 (0x00000001) |

### refenerces
* https://stackoverflow.com/questions/9115436/performance-counter-by-process-id-instead-of-name
* https://docs.microsoft.com/en-us/dotnet/framework/debug-trace-profile/performance-counters-and-in-process-side-by-side-applications

## Configuration:

```toml
# Add username win_proc metric
[[processors.proc_username]]
```
This plugin depends on win_proc and win_dotnet_* measurements
```toml
[[inputs.win_perf_counters.object]]
  # Process metrics, in this case for IIS only
  ObjectName = "Process"
  Counters = ["% Processor Time","Thread Count","Working Set - Private"]
  Instances = ["w3wp*"]
  Measurement = "win_proc"

[[inputs.win_perf_counters.object]]
  # .NET CLR Memory, in this case for IIS only
  ObjectName = ".NET CLR Memory"
  Counters = ["% Time in GC","# Bytes in all Heaps"]
  Instances = ["w3wp*"]
  Measurement = "win_dotnet_mem"
```

## Tags:

username="username"
