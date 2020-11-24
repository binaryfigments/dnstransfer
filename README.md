# dnstransfer

This tool (`dnstransfer`) is a commandline tool that can check multiple domain names for zone transfers. (axfx)

# Logic

1. Check per domain nameservers (like `dig NS zonetransfer.me`)
2. Do AXFR lookup per domain per nameserver (like `dig AXFR zonetransfer.me @nsztm1.digi.ninja.`)

# Help

```shell
user@laptop:~$ ./dnstransfer check -h
This command checks a list of domain names for zone transfers.

Usage:
  dnstransfer check [flags]

Flags:
      --file string   File with domain names. (default "domains.txt")
  -h, --help          help for check

Global Flags:
      --config string       config file (default is $HOME/.dnstransfer.yaml)
  -d, --debug               all output for debugging
  -j, --json                JSON output
      --logfile string      Output to logfile.
      --nameserver string   What nameserver to use. (default "8.8.8.8")
```

# Example

Check domain names from file on AXFR's.


```shell
sebastian@vpn:~$ ./dnstransfer check --file domains.txt 
ERRO[0000] Zone can be transfered!                       domain=alblasserwaard-vijfheerenlanden.nl error="<nil>" nameserver=dns1.qdc.nl. transferable=true
ERRO[0001] Zone can be transfered!                       domain=parkstad-limburg.nl error="<nil>" nameserver=dns2.qdc.nl. transferable=true
ERRO[0001] Zone can be transfered!                       domain=advlimburg.nl error="<nil>" nameserver=nszero2.axc.nl. transferable=true
ERRO[0001] Zone can be transfered!                       domain=tresoar.nl error="<nil>" nameserver=nszero2.axc.nl. transferable=true
ERRO[0001] Zone can be transfered!                       domain=nef.nl error="<nil>" nameserver=nszero1.axc.nl. transferable=true
ERRO[0001] Zone can be transfered!                       domain=advlimburg.nl error="<nil>" nameserver=nszero1.axc.nl. transferable=true
ERRO[0001] Zone can be transfered!                       domain=ggdflevoland.nl error="<nil>" nameserver=dns1.qdc.nl. transferable=true
ERRO[0001] Zone can be transfered!                       domain=parkstad-limburg.nl error="<nil>" nameserver=dns1.qdc.nl. transferable=true
ERRO[0001] Zone can be transfered!                       domain=ggdflevoland.nl error="<nil>" nameserver=dns2.qdc.nl. transferable=true
ERRO[0001] Zone can be transfered!                       domain=tresoar.nl error="<nil>" nameserver=nszero1.axc.nl. transferable=true
ERRO[0002] Zone can be transfered!                       domain=visitaties.nl error="<nil>" nameserver=ns3.qweb.co. transferable=true
ERRO[0002] Zone can be transfered!                       domain=randstedelijke-rekenkamer.nl error="<nil>" nameserver=nszero1.axc.nl. transferable=true
ERRO[0002] Zone can be transfered!                       domain=veere.nl error="<nil>" nameserver=ns2.qweb.nu. transferable=true
ERRO[0002] Zone can be transfered!                       domain=visitaties.nl error="<nil>" nameserver=ns2.qweb.nl. transferable=true
ERRO[0002] Zone can be transfered!                       domain=randstedelijke-rekenkamer.nl error="<nil>" nameserver=nszero2.axc.nl. transferable=true
ERRO[0002] Zone can be transfered!                       domain=visitaties.nl error="<nil>" nameserver=ns1.qweb.nl. transferable=true
ERRO[0002] Zone can be transfered!                       domain=veere.nl error="<nil>" nameserver=ns3.qweb.co. transferable=true
ERRO[0002] Zone can be transfered!                       domain=veere.nl error="<nil>" nameserver=ns1.qweb.net. transferable=true
```

# Notes

```shell
cobra init --author "Sebastian Broekhoven @binaryfigments" --license MIT --pkg-name github.com/binaryfigments/dnstransfer
```

Added "check" command

```shell
cobra add check -t github.com/binaryfigments/dnstransfer
```