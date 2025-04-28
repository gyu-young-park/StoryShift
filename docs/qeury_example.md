1. Get posts with limit

- request
```sh
http://localhost:8080/v1/velog/posts?name=chappi&limit=3&url_slog=6d7991cc-0cb3-4dfa-b212-8de9f2476663
```

- response
```json
{
  "posts": [
    {
      "id": "4023bf7e-df1c-4288-9e4f-a37983406912",
      "title": "eBPF를 배워보자 3일차 - eBPF Program 해부",
      "created_at": "2025-04-22T12:36:08.446Z",
      "updated_at": "2025-04-26T10:10:15.629Z",
      "short_description": "이전에는 bcc를 사용해서 eBPF 사용해보았는데, 이제는 c언어를 직접 사용하여 bcc가 어떻게 동작했는 지 알아보도록 하자.c또는 Rust source code는 eBPF bytecode로 컴파일된다. 이 eBPF bytecode는 JIT compile되거나 int",
      "thumnail": "https://velog.velcdn.com/images/chappi/post/42f30731-7b86-4713-a429-3acc63d288a1/image.png",
      "url_slog": "eBPF를-배워보자-3일차-eBPF-Program-해부",
      "tags": [
        "ebpf",
        "linux"
      ]
    },
    {
      "id": "edcce6f5-581c-464f-a202-21ba6792ae62",
      "title": "eBPF를 배워보자 2일차 - eBPF \"Hello World\"",
      "created_at": "2025-04-22T12:33:38.459Z",
      "updated_at": "2025-04-24T07:21:18.353Z",
      "short_description": "만약 필요한 toolchain들이 없다면 설치해주도록 하자. ebpf를 실행하기 위해서는 libbpf가 필요하다. libbpf는 BPF프로그램을 compile하고 load하는 데 사용되는 C라이브러리이다. 즉, BPF 프로그램의 실행에 초점을 맞춘 도구라고 볼 수 있는",
      "thumnail": "https://velog.velcdn.com/images/chappi/post/c7797b46-fbed-457b-a74c-8e2e0963ff5d/image.png",
      "url_slog": "eBPF를-배워보자-2일차-eBPF-Hello-World",
      "tags": [
        "ebpf",
        "linux"
      ]
    },
    {
      "id": "d0295770-ea06-48a8-bc39-2f8f02182f4c",
      "title": "eBPF를 배워보자 1일차 - eBPF란?",
      "created_at": "2025-04-22T12:29:46.672Z",
      "updated_at": "2025-04-26T05:16:49.245Z",
      "short_description": "What is eBPF, and why is it important\neBPF는 custom code를 작성하여 kernel에 동적으로 적재하여 kernel의 동작을 변경할 수 있다. 이를 통해서 좋은 성능의 네트워킹, observability, security tool",
      "thumnail": "https://velog.velcdn.com/images/chappi/post/9eeb7980-69a2-43a8-a578-dad3b7e620f7/image.png",
      "url_slog": "eBPF",
      "tags": [
        "ebpf",
        "linux"
      ]
    }
  ]
}
```

2. Get post information in detail
- request
```sh
http://localhost:8080/v1/velog/post?name=chappi&url_slog=eBPF
```

- response
```json
{
  "post": {
    "id": "d0295770-ea06-48a8-bc39-2f8f02182f4c",
    "title": "eBPF를 배워보자 1일차 - eBPF란?",
    "created_at": "2025-04-22T12:29:46.672Z",
    "updated_at": "2025-04-26T05:16:49.245Z",
    "body": "# What is eBPF, and why is it important\neBPF는 custom code를 작성하여 ... (생략)..."
  }
}
```

3. Download post

- request
```sh
http://localhost:8080/v1/velog/post/download?name=chappi&url_slog=eBPF
```

- response
```
eBPF를+배워보자+1일차+-+eBPF란_.md
```