# LinuxCNC UI Experimentation

This repository contains a C-Wrapper and Go adapter for the LinuxCNC API.

## Testing

Building and testing the code in this repository requires to compile LinuxCNC from source.
Plus Go 1.17 installed.

```sh
# make sure to be in the rip environment
. /home/nschieder/linuxcnc/scripts/rip-environment

# Dump C struct data obtained via Go to stdout.
# Also builds the liblinuxcncadapter.so shared library.
go run ./cmd/mage run:datadump
```

Example output when LinuxCNC is running in the background using `axis_mm.ini`.
```
global:
 main._Ctype_struct_stat_Global{estop:0, cycleTime:0.001, trajectoryPlannerEnabled:1, file:(*main._Ctype_char)(0xd40b00), inpos:true, motionPaused:false, mist:0, flood:0, toolInSpindle:0, pocketPrepped:-1, numberOfSpindles:1, numberOfJoints:3, _:[4]uint8{0x0, 0x0, 0x0, 0x0}}

joint #0 main._Ctype_struct_stat_Joint{backlash:0, enabled:1, fault:0, ferrorCurrent:3.4858793911303476e-17, ferrorHighMark:1.1102230246251565e-16, homed:1, homing:0, inpos:1, input:2.5691630378688046e-08, jointType:1, maxFerror:1.27, maxHardLimit:0, maxPositionLimit:254, maxSoftLimit:0, minFerror:0.254, minHardLimit:0, minPositionLimit:-254, minSoftLimit:0, output:2.569163041354684e-08, overrideLimits:0, units:1, velocity:0}
joint #1 main._Ctype_struct_stat_Joint{backlash:0, enabled:1, fault:0, ferrorCurrent:4.508184630995244e-17, ferrorHighMark:1.1102230246251565e-16, homed:1, homing:0, inpos:1, input:2.845020996566916e-09, jointType:1, maxFerror:1.27, maxHardLimit:0, maxPositionLimit:254, maxSoftLimit:0, minFerror:0.254, minHardLimit:0, minPositionLimit:-254, minSoftLimit:0, output:2.8450210416487623e-09, overrideLimits:0, units:1, velocity:0}
joint #2 main._Ctype_struct_stat_Joint{backlash:0, enabled:1, fault:0, ferrorCurrent:-7.655051790481267e-16, ferrorHighMark:1.7763568394002505e-15, homed:1, homing:0, inpos:1, input:2.9193738981803108e-08, jointType:1, maxFerror:1.27, maxHardLimit:0, maxPositionLimit:101.6, maxSoftLimit:0, minFerror:0.254, minHardLimit:0, minPositionLimit:-50.8, minSoftLimit:0, output:2.919373821629793e-08, overrideLimits:0, units:1, velocity:0}

spindle #0 main._Ctype_struct_stat_Spindle{brake:1, direction:0, enabled:0, override:1, overrideEnabled:true, speed:0}
```
