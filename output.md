# chapter 3. eBPF Program 해부
이전에는 `bcc`를 사용해서 eBPF 사용해보았는데, 이제는 c언어를 직접 사용하여 `bcc`가 어떻게 동작했는 지 알아보도록 하자.

c또는 Rust source code는 eBPF bytecode로 컴파일된다. 이 eBPF bytecode는 JIT compile되거나 interpreted되어 native machine code 명령어로 변환된다. 다음의 그림을 참고하자.  
![](https://velog.velcdn.com/images/chappi/post/42f30731-7b86-4713-a429-3acc63d288a1/image.png)

eBPF program은 eBPF bytecode 명령어 셋으로 assembly로 programming할 수 있지만, 사람이 읽을 수 있고, programming하기 좋은 c나 rust와 같은 언어로 먼저 작성하고 bytecode로 만들어 실행하는 것이 좋다. 

결과적으로 eBPF bytecode는 kernel의 eBPF virtual machine안에서 실행이 되는 것이다. 

## 1. eBPF Virtual Machine
eBPF virtual machine은 다른 virtual machine과 같이 computer의 software 구현체이다. virtual machine은 eBPF bytecode 명령어 형식으로 program을 받아들고 bytecode를 CPU에 해당하는 native machine 명령어로 변환해준다. 

eBPF 초기에 bytecode 명령어들은 kernel안에서 interpreted되었다. 즉, eBPF program이 실행될 때마다 kernel은 명령어를 평가하고 bytecode를 machine code로 변환해 실행했다는 것이다. 이후에 eBPF interpreter의 성능 이슈와 취약성을 피하기 위해 interpreter를 JIT compilation으로 변경하였다. 이 덕분에 bytecode를 한번만 컴파일해 native machine 명령어로 변환하게되었고, kernel에 program을 한 번만 loading시켜 code를 구동시킬 수 있게되었다. 

eBPF bytecode는 명령어 셋을 포함하며 이러한 명령어들은 (virtual) eBPF register에서 동작한다. eBPF 명령어 셋과 register model은 CPU 아키첵쳐에 맞게 동작하도록 설계되었기 때문에, bytecode를 machine code로 컴파일하거나 인터프리팅하는 과정은 매우 직관적이다. 

eBPF virtual machine은 10개의 general-purpose register를 사용하여 0~9까지 번호로 되어있다. 추가적으로 register 10은 stack frame pointer로 사용되어 read만 가능하고 write는 불가능하다. BPF pogram이 실행되면 value의 state를 기록하기 위해 register에 value를 기록한다. 

register는 `BPF_REG_0` ~ `BPF_REG_10`까지 `bpf.h` 파일에 적혀있다.  

실행되기전에 eBPF program에 대한 context argument는 Register 1번에 로드되고, 함수의 return value는 Register 0번에 저장된다.

eBPF code로부터 함수를 실행하기 이전에 함수에 대한 argument들은 Register 1~5까지 위치하게 되는 것이다. 

`linux/bpf.h` header file은 `bpf_insn`라는 구조체를 정의하였는데 이는 BPF 명령어를 나타낸다.
```c
struct bpf_insn {
    __u8 code;          /* opcode */
    __u8 dst_reg:4;     /* dest register */
    __u8 src_reg:4;     /* source register */
    __s16 off;       /* signed offset */
    __s32 imm;       /* signed immediate constant */
};

\n
```
`code`는 opcode를 말한다. 각 명령어는 opcode를 가진다. 이는 명령어가 실행하는 operation을 표현하는 것이다. 가령 register에 value를 저한다거나 program의 다른 명령어로 jump하는 경우가 있다. 

`bpf_insn` 구조체는 64-bit(8 byte) long이다. 하지만 명령어는 8byte이상을 필요로 할 때가 있다. 이럴 때에는 wide 명령어 encoding을 사용하여 16 bytes long을 표현한다. 이에 대해서는 추후에 알아보도록 하자.

kernel에 로드될 때 bytecode가 로드될 때, eBPF program의 bytecode는 일련의 eBPF 명령어 모음인 `bpf_insn` 구조체들로 표현된다. verifier는 eBPF 명령어들으 분석하여 해당 code가 안전하게 동작할 지 검사하고 보장해준다. 

대부분의 opcode는 다음의 카테고리로 빠져나간다.
1. register에 value를 적재하기
2. register의 값을 memory에 저장하기
3. 산술 연산 수행하기
4. 특정 조건이 만족하면 다른 명령어로 jump하기

이제 C언어로 eBPF code를 만들어서 사용해보도록 하자.

## 2. Network Interface에서 eBPF "Hello World"
network packet이 도착하면 실행되는 eBPF program을 만들어보도록 하자.

`hello.bpf.c` file을 만들어서 network packet이 들어오면 counter와 함께 `Hello World`를 출력하도록 하는 eBPF program을 만들었다. file이름 가운에 `bpf`를 넣는 것은 하나의 관례로 생각하면 된다.

- hello.bpf.c
```c
#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>

int counter = 0;

SEC("xdp")
int hello(struct xdp_md *ctx) {
    bpf_printk("Hello World %d", counter);
    counter++; 
    return XDP_PASS;
}

char LICENSE[] SEC("license") = "Dual BSD/GPL";
```
먼저 코드 설명 뒤에 실행하는 방법에 대해서 알아보도록 하자.

1. `SEC()` macro는 `xdp`라고 불리는 section을 정의하는 것이다. 추후에 알아보도록 하고 지금은 eBPF program의 XDP(eXpress Data Path)를 정의하고 있다는 것만들 알아두도록 하자.
2. `hello` eBPF program 이름을 볼 수 있는데, 재밌는 것은 program 이름이 function 이름이라는 것이다. 따라서 program을 `hello`라고 부르는 것이다. helper function으로 `bpf_printk`를 사용하여 log를 찍고 `XDP_PASS`반환한다. 이는 network packet을 정상적으로 처리하라는 것을 kernel에 알려주는 것이다. 
3. `char LICENSE[] SEC("license") = "Dual BSD/GPL";`에서도 `SEC()`가 사용되는 것을 볼 수 있다. 이는 해당 eBPF program을 실행할 라이센스를 알려주는 code이다. 어떤 helper function은 `GPL only`라고 적혀있는데 이는 오직 GPL-compatible license에서만 가능하다는 것이다. `GPL only`로 적힌 helper function에 대해서는 GPL-compatible lince를 가지는 것으로 선언해야 BPF code에서 해당 helper function을 사용할 수 있다. 이에 대해서는 추후에 더 자세히 알아보도록 하자.

위의 예제는 eBPF program을 network interface의 XDP hook point에 attach한 code이다. network packet이 physical 또는 virtual network interface에 들어올 때 XDP event가 발생하고 해당 eBPF program이 실행된다고 생각하면 된다. 


이제 c언어로 된 ebpf 빌드 방법이다. 참고로 `clang`은 version 12이상으로 실행하는 것이 좋다고 한다.
```sh
clang -target bpf -I/usr/include/$(uname -m)-linux-gnu -g -O2 -c ./hello.bpf.c -o ./hello.bpf.o
```
위 code는 compile되어 eBPF virtual machine이 이해할 수 있는 eBPF bytecode로 변환되어야 한다. 이때 `-target bpf`를 사용해야한다. 위의 명령어를 실행하면 `hello.bpf.o` object파일을 `hello.bpf.c` file로 부터 생성해내는데 `-g`는 optional인데, 아래에서 알아보겠지만 data를 좀 더 인간이 보기좋게 표현화해준다. 만약 이 옵션을 설정해주지 않으면 바이트 값으로 디버깅해야한다. 

아래의 에러 발생 시에
```sh
In file included from hello.bpf.c:1:
In file included from /usr/include/linux/bpf.h:11:
/usr/include/linux/types.h:5:10: fatal error: 'asm/types.h' file not found
#include <asm/types.h>
```

다음의 설치하면 해결할 수 있다.
```sh
sudo apt-get install -y gcc-multilib
```

이제 `eBPF` object file을 분석해보도록 하자. 가장 흔히 사용될 수 있는 명령어가 `file`이다.
```sh
file hello.bpf.o 
hello.bpf.o: ELF 64-bit LSB relocatable, eBPF, version 1 (SYSV), with debug_info, not stripped
```
`ELF`는 executable and linkable format file이라는 이고, eBPF code이며 64bit-platform에 LSB(least significant bit) 아키텍처를 사용하고 있다는 것을 알려주고 있다. 

`llvm-objdump`를 사용하여 eBPF 명령어를 더 자세히 분석할 수 있다.
```sh
llvm-objdump -S hello.bpf.o 

hello.bpf.o:    file format ELF64-BPF


Disassembly of section xdp:

0000000000000000 hello:
;     bpf_printk("Hello World %d", counter);
       0:       18 06 00 00 00 00 00 00 00 00 00 00 00 00 00 00 r6 = 0 ll
       2:       61 63 00 00 00 00 00 00 r3 = *(u32 *)(r6 + 0)
       3:       18 01 00 00 00 00 00 00 00 00 00 00 00 00 00 00 r1 = 0 ll
       5:       b7 02 00 00 0f 00 00 00 r2 = 15
       6:       85 00 00 00 06 00 00 00 call 6
;     counter++; 
       7:       61 61 00 00 00 00 00 00 r1 = *(u32 *)(r6 + 0)
       8:       07 01 00 00 01 00 00 00 r1 += 1
       9:       63 16 00 00 00 00 00 00 *(u32 *)(r6 + 0) = r1
;     return XDP_PASS;
      10:       b7 00 00 00 02 00 00 00 r0 = 2
      11:       95 00 00 00 00 00 00 00 exit
```
`disassembly`에 대해서 잘 모르지만 위의 코드를 명시적으로 이해하는 것은 크게 어렵지 않다.

`xdp`로 label된 section은 c code에서 `SEC()`에 해당하는 부분이고 해당 section의 함수는 `hello`인 것이다. `bpf_printk`에 해당하는 eBPF bytecode 명령어는 5개의 line으로 된 것이다. `coutner`는 3개의 bytecode 명령어 line으로 `return XDP_PASS`는 두 개의 bytecode 명령어 line으로 되어있는 것이다.

## 3. kernel에 program 로딩 및 로딩된 program 분석하기
이번에는 `bpftool`을 사용해볼 것이다. `bpftool`은 program을 kernel에 load시켜 동작을 확인하고 debug한다. 물론 program을 programmatically하게 load하는 방식도 있지만 이에 대해서는 추후에 알아보도록 하자.

다음의 예제는 `bpftool`을 사용해서 program을 kernel에 load하는 예제이다. 
```sh
bpftool prog load hello.bpf.o /sys/fs/bpf/hello
```

위는 `hello.bpf.o` object file을 load하고 `/sys/fs/bpf/hello`에 고정시키도록 한 것이다. load에 성공하면 어떠한 응답도 없다. `ls`를 통해 `/sys/fs/bpf/hello`에서 결과를 확인할 수 있다.
```sh
ls /sys/fs/bpf/
hello
```
eBPF program이 성공적으로 로딩된 것을 알 수 있다. 이제 `bpftool`을 사용해서 kernel안에서의 bpf program과 그 status에 대해서 확인해보도록 하자.

`bpftool`은 kernel에 loading된 모든 program을 확인해준다.
```sh
bpftool prog list
...
157: xdp  name hello  tag d35b94b4c0c10efb  gpl
        loaded_at 2024-01-30T15:37:15+0900  uid 0
        xlated 96B  jited 63B  memlock 4096B  map_ids 43,44
        btf_id 239
```

`hello` eBPF program은 `157` ID를 배정받았다. 이 ID는 kernel에 로딩되면 배정받는 ID값으로 이를 기반으로 program에 대한 정보를 얻어낼 수있다.
```sh
bpftool prog show id 157 --pretty
{
    "id": 157,
    "type": "xdp",
    "name": "hello",
    "tag": "d35b94b4c0c10efb",
    "gpl_compatible": true,
    "loaded_at": 1706596635,
    "uid": 0,
    "orphaned": false,
    "bytes_xlated": 96,
    "jited": true,
    "bytes_jited": 63,
    "bytes_memlock": 4096,
    "map_ids": [43,44
    ],
    "btf_id": 239
}
```
file name과 같이 명확하게 그 값의 의미를 알 수 있다. 

1. `id`: 해당 eBPF id는 157이다.
2. `type`: 이 eBPF program이 XDP event를 사용해 network interface에 attach될 수 있다는 것을 나타낸다.
3. `name`: 해당 program의 이름이다. 이 이름은 source code의 함수 이름과 같다.
4. `tag`: program의 또다른 식별자로 나중에 더 자세히 다루어보자.
5. `gpl_compatible`: 해당 program이 GPL-compatible license를 가지고 있다는 것이다.
6. `loaded_at`: timestamp로 언제 해당 program이 로딩되었는 지를 나타낸다.
7. `uid`: user id는 해당 program을 loading한 program으로 root인 0번이다.
8. `bytes_jited`: eBPF program을 machine code로 compile한 결과가 63 byte라는 것이다. 
9. `bytes_memlock`: 해당 program은 4096 byte memory를 보존하고 있다는 것이다.
10. `map_ids`: 해당 program이 165, 166 ID를 가진 BPF map을 참조하고 있다는 것이다. source code에서는 명시적으로 BPF map을 사용하지 않았지만 eBPF program이 global data에 접근하기 위해서 임의적으로 만든 것이다. 이에 대해서는 추후에 자세히 알아보도록 하자.
11. `btf_id`: `btf_id`는 해당 program을 위한 BTF 정보 block이 있다는 것을 나타낸다. 이 정보는 `-g` flag로 comple할 때만 object file에 포함된다.

`tag`는 SHA(secure hashing algorithm)으로 program 명령어들의 합이다. 이는 program에 대한 또 다른 식별자로 쓰인다. program의 ID는 load, unload가 반복되면 계속 바뀌지만 `tag`는 계속 동일하게 남아있다. `bpftool`은 BPF program을 ID, name, tag 또는 path를 통해 참조한다. 따라서 다음과 같이 program을 참조할 수 있다.

- `bpftool prog show id 540`
- `bpftool prog show name hello`
- `bpftool prog show tag d35b94b4c0c10efb`
- `bpftool prog show pinned /sys/fs/bpf/hello`

여러 program에 같은 name 또는 tag를 달 수는 있지만 id와 고정한 path는 반드시 unique해야한다.

`bytes_xlated` file는 eBPF code가 verifier 통과하고 얼마만큼의 byte로 구성되어있는 지를 나타낸다. `bpftool`을 사용해서 확인해보도록 하자.
```sh
bpftool prog dump xlated name hello 
int hello(struct xdp_md * ctx):
; bpf_printk("Hello World %d", counter);
   0: (18) r6 = map[id:43][0]+0
   2: (61) r3 = *(u32 *)(r6 +0)
   3: (18) r1 = map[id:44][0]+0
   5: (b7) r2 = 15
   6: (85) call bpf_trace_printk#-64880
; counter++; 
   7: (61) r1 = *(u32 *)(r6 +0)
   8: (07) r1 += 1
   9: (63) *(u32 *)(r6 +0) = r1
; return XDP_PASS;
  10: (b7) r0 = 2
  11: (95) exit
```
위 code는 `llvm-objdump`의 결과물과 매우 유사한 disaasemly값이다. 

eBPF bytecode는 JIT compiler에 의해서 CPU아키텍처에 맞는 machine code에 변환된다. `bytes_jited` field는 program이 변환된 후에 machine code로 108 bytes long으로 구성되어 있다는 것을 나타낸다.

`bpftool`을 통해서 assembly형식의 jit-compiled code를 dump를 만들어줄 수 있다.
```sh
bpftool prog dump xlated name hello 
int hello(struct xdp_md * ctx):
; bpf_printk("Hello World %d", counter);
   0: (18) r6 = map[id:43][0]+0
   2: (61) r3 = *(u32 *)(r6 +0)
   3: (18) r1 = map[id:44][0]+0
   5: (b7) r2 = 15
   6: (85) call bpf_trace_printk#-64880
; counter++; 
   7: (61) r1 = *(u32 *)(r6 +0)
   8: (07) r1 += 1
   9: (63) *(u32 *)(r6 +0) = r1
; return XDP_PASS;
  10: (b7) r0 = 2
  11: (95) exit


bpftool prog dump jited name hello 
int hello(struct xdp_md * ctx):
bpf_prog_d35b94b4c0c10efb_hello:
; bpf_printk("Hello World %d", counter);
   0:   nopl    (%rax,%rax)
   5:   nop
   7:   pushq   %rbp
   8:   movq    %rsp, %rbp
   b:   pushq   %rbx
   c:   movabsq $-92458832723968, %rbx
  16:   movl    (%rbx), %edx
  19:   movabsq $-110869714558704, %rdi
  23:   movl    $15, %esi
  28:   callq   0xffffffffce247d40
; counter++; 
  2d:   movl    (%rbx), %edi
  30:   addq    $1, %rdi
  34:   movl    %edi, (%rbx)
; return XDP_PASS;
  37:   movl    $2, %eax
  3c:   popq    %rbx
  3d:   leave
  3e:   retq
```
assembly language를 하나하나 이해할 필요는 없고, 우리의 eBPF code가 다음과 같은 assembly언어로 JIT-compiler에 의해 변환된다는 사실에 집중하도록 하자.

## 4. Attaching to an Event
program type은 program이 attach된 event의 type과 매치해야한다. 우리의 예제에서는 XDP program으로 `bpftool`을 사용해서 network interface에 대한 `XDP` event에 program을 attach할 수 있다.
```sh
bpftool net attach xdp id 157 dev eth0
```
eBPF program id 157을 network interface인 `eth0`에 attach시켰다.

`bpftool`을 사용해서 eBPF program에 연결된 network interface를 볼 수 있다.
```sh
bpftool net list
xdp:
eth0(2) driver id 157

tc:

flow_dissector:
```
해당 eBPF program은 eth0 interface에서 발생하는 `XDP` event에 attach되었다. 

우리의 `hello` BPF program은 network packet이 들어올때마다 trace output을 출력할 것이다. `cat /sys/kernel/debug/tracing/trace_pipe`를 사용해도 되지만, `bpftool prog tracelog`를 사용해도 된다. 
```sh
bpftool prog tracelog
...
<idle>-0       [003] d.s.. 655370.944105: bpf_trace_printk: Hello World 4531
<idle>-0       [003] d.s.. 655370.944587: bpf_trace_printk: Hello World 4532
<idle>-0       [003] d.s.. 655370.944896: bpf_trace_printk: Hello World 4533
```
`XDP` event들이 network packet이 도착할 때마다 발생하여 우리의 eBPF program이 실행되는 것을 볼 수 있다. 시간이 지나면 알아서 network interface로의 연결이 끝어지므로 참고하도록 하자.

`4531`, `4532`와 같이 counter가 증가하는 것을 볼 수 있다. eBPF에서는 global variable에 대해서 어떻게 관리하고 다루는 지 알아보도록 하자.

eBPF map은 data 구조로 eBPF program으로부터 또는 user space으로부터 접근할 수 있다. map은 다른 program에서도 접근할 수 있기 때문에 map은 value의 state를 유지시키고 다음 실행 program으로 넘겨야한다. 여러 program들이 같은 map에 접근할 수 있는 것이다. 이러한 특징 때문에 map의 의미는 global variable로 변형되었다. 

`bpftool`을 사용하면 kernel에 로딩된 map을 보여준다. 이전에 `43`과 `44` map을 사용한다고 나와있었는데, 확인해보도록 하자.
```sh
bpftool map list
43: array  name hello.bss  flags 0x400
        key 4B  value 4B  max_entries 1  memlock 4096B
        btf_id 239
44: array  name hello.rodata  flags 0x80
        key 4B  value 15B  max_entries 1  memlock 4096B
        btf_id 239  frozen
```
C source program으로부터 컴파일된 object file안의 bss영역은 global variable들을 가지고 있다는 것을 알 수 있다. `bpftool`을 사용하여 그 내부를 들여다 볼 수 있다.
```sh
bpftool map dump name hello.bss 
[{
        "value": {
            ".bss": [{
                    "counter": 10739
                }
            ]
        }
    }
]
```
`bpftool map dump id 43`으로도 가능하다. 우리의 `counter`값이 `.bss`영역안에 있는 것을 확인할 수 있다. 

map은 또한 static data를 가지기 위해 사용되는데 이는 `hello.rodata` map을 확인하면 된다. `hello.rodata`는 read-only data라는 의미로 `hello`에 관한 metadata가 들어있다고 보면 된다.
```sh
bpftool map dump name hello.rodata
[{
        "value": {
            ".rodata": [{
                    "hello.____fmt": "Hello World %d"
                }
            ]
        }
    }
]
```

참고로 이렇게 data가 이쁘게 보이는 이유는 우리가 빌드할 때 `-g`옵션을 썻기 때문이다. 만약 `-g`옵션을 스지않으면 굉장히 보기 힘든 바이트값들이 나온다. 

## 5. eBPF program 종료하기
network interface로 부터 eBPF program을 detach시키고 싶다면 다음의 명령어를 쓰도록 하자.
```sh
bpftool net detach xdp dev eth0
```
다음의 명령어로 network interface에 대한 연결을 종료할 수 있다. 어떠한 output도 없다면 성공적으로 명령어를 실행한 것이다. 

```sh
bpftool net list
xdp:

tc:

flow_dissector:
...
```
다음과 같이 비어있다면 성공이다. 

그러나 program은 여전히 kernel에 로딩되어있다. 
```sh
bpftool prog show name hello 
157: xdp  name hello  tag d35b94b4c0c10efb  gpl
        loaded_at 2024-01-30T15:37:15+0900  uid 0
        xlated 96B  jited 63B  memlock 4096B  map_ids 43,44
        btf_id 239
```

고정된 pseudofile을 삭제함으로서 bpf program을 종료시킬 수 있다.
```sh
rm /sys/fs/bpf/hello
bpftool prog show name hello 
```
별다른 output이 없다면 성공이다. 

## 6. BPF to BPF calls
eBPF program에서 다른 함수를 실행하기 위해서는 `tail calls`를 사용해야한다고 했다. 그러나 이제는(linux kerne 4.16과 LLVM 6.0이후로) eBPF program안에서도 다른 함수를 호출할 수 있는 기능이 있다고 언급했었다. 가능한 지 확인해보도록 하자.

- hello-func.bpf.c
```c
#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>

static __attribute((noinline)) int get_opcode(struct bpf_raw_tracepoint_args *ctx) {
    return ctx->args[1];
}

SEC("raw_tp/")
int hello(struct bpf_raw_tracepoint_args *ctx) {
    int opcode = get_opcode(ctx);
    bpf_printk("Syscall: %d", opcode);
    return 0;
}

char LICENSE[] SEC("license") = "Dual BSD/GPL";
```
`get_opcode`함수는 `args[1]`을 추출해내는 것이 전부인데, `args[1]`이 바로 실행된 syscall의 `opcode`이다. 

`__attribute((noinline))`은 해당 function을 `inline`으로 만들지 않도록 하기 위해서 선언한 지시문이다. `eBPF` function은 다음의 함수를 아래와 같이 쓸 수 있다.
```c
SEC("raw_tp/")
int hello(struct bpf_raw_tracepoint_args *ctx) {
    int opcode = get_opcode(ctx);
    bpf_printk("Syscall: %d", opcode);
    return 0;
}
```
놀랍게도 별다른 작업없이 사용할 수 있다.

eBPF source code를 object file로 컴파일해보도록 하자.
```sh
clang -target bpf -I/usr/include/$(uname -m)-linux-gnu -g -O2 -c ./hello-func.bpf.c -o ./hello-func.bpf.o
```

eBPF program을 실행시켜보도록 하자.
```sh
bpftool prog load hello-func.bpf.o /sys/fs/bpf/hello
bpftool prog list name hello 
244: raw_tracepoint  name hello  tag 3d9eb0c23d4ab186  gpl
        loaded_at 2024-01-30T17:59:58+0900  uid 0
        xlated 80B  jited 60B  memlock 4096B  map_ids 88
        btf_id 305
```
`get_opcode()`함수를 보기위해서 eBPF bytecode를 검사해보도록 하자.

```sh
bpftool prog dump xlated name hello 
int hello(struct bpf_raw_tracepoint_args * ctx):
; int opcode = get_opcode(ctx);
   0: (85) call pc+7#bpf_prog_cbacc90865b1b9a5_get_opcode
; bpf_printk("Syscall: %d", opcode);
   1: (18) r1 = map[id:88][0]+0
   3: (b7) r2 = 12
   4: (bf) r3 = r0
   5: (85) call bpf_trace_printk#-64880
; return 0;
   6: (b7) r0 = 0
   7: (95) exit
int get_opcode(struct bpf_raw_tracepoint_args * ctx):
; return ctx->args[1];
   8: (79) r0 = *(u64 *)(r1 +8)
; return ctx->args[1];
   9: (95) exit
```
`hello` eBPF program이 `get_opcode()`함수를 호출하는 것을 볼 수 있다. `get_opcode()` byte code도 있는 것을 확인할 수 있다. 단, stack size가 512 byte로 한정되어 있기 때문에 BPF에서 BPF 함수를 깊게 호출할 수 없다. 

