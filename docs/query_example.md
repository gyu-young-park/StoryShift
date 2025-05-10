1. Get posts with limit

- request
```sh
http://localhost:9596/v1/velog/chappi/posts?count=10
```

- response
```json
{
  "posts": [
    {
      "id": "4023bf7e-df1c-4288-9e4f-a37983406912",
      "title": "eBPF를 배워보자 3일차 - eBPF Program 해부",
      "created_at": "2025-04-22T12:36:08.446Z",
      "updated_at": "2025-04-30T06:14:34.31Z",
      "short_description": "이전에는 bcc를 사용해서 eBPF 사용해보았는데, 이제는 c언어를 직접 사용하여 bcc가 어떻게 동작했는 지 알아보도록 하자.c또는 Rust source code는 eBPF bytecode로 컴파일된다. 이 eBPF bytecode는 JIT compile되거나 int",
      "thumnail": "https://velog.velcdn.com/images/chappi/post/42f30731-7b86-4713-a429-3acc63d288a1/image.png",
      "url_slug": "eBPF를-배워보자-3일차-eBPF-Program-해부",
      "tags": [
        "ebpf",
        "linux"
      ]
    },
    {
      "id": "edcce6f5-581c-464f-a202-21ba6792ae62",
      "title": "eBPF를 배워보자 2일차 - eBPF \"Hello World\"",
      "created_at": "2025-04-22T12:33:38.459Z",
      "updated_at": "2025-04-30T02:26:14.255Z",
      "short_description": "만약 필요한 toolchain들이 없다면 설치해주도록 하자. ebpf를 실행하기 위해서는 libbpf가 필요하다. libbpf는 BPF프로그램을 compile하고 load하는 데 사용되는 C라이브러리이다. 즉, BPF 프로그램의 실행에 초점을 맞춘 도구라고 볼 수 있는",
      "thumnail": "https://velog.velcdn.com/images/chappi/post/c7797b46-fbed-457b-a74c-8e2e0963ff5d/image.png",
      "url_slug": "eBPF를-배워보자-2일차-eBPF-Hello-World",
      "tags": [
        "ebpf",
        "linux"
      ]
    }
  ]
}
```

2. Get posts with limit and post_id

- request
```sh
http://localhost:9596/v1/velog/chappi/posts?count=2&post_id=edcce6f5-581c-464f-a202-21ba6792ae62
```

- response
```sh
{
  "posts": [
    {
      "id": "d0295770-ea06-48a8-bc39-2f8f02182f4c",
      "title": "eBPF를 배워보자 1일차 - eBPF란?",
      "created_at": "2025-04-22T12:29:46.672Z",
      "updated_at": "2025-04-30T05:01:27.35Z",
      "short_description": "What is eBPF, and why is it important\neBPF는 custom code를 작성하여 kernel에 동적으로 적재하여 kernel의 동작을 변경할 수 있다. 이를 통해서 좋은 성능의 네트워킹, observability, security tool",
      "thumnail": "https://velog.velcdn.com/images/chappi/post/9eeb7980-69a2-43a8-a578-dad3b7e620f7/image.png",
      "url_slug": "eBPF",
      "tags": [
        "ebpf",
        "linux"
      ]
    },
    {
      "id": "56625df5-632d-4ad9-b6fb-8c67b1624704",
      "title": "SQL 재활 훈련 9일차 - View와 Having",
      "created_at": "2025-04-15T04:06:17.341Z",
      "updated_at": "2025-04-24T13:32:50.551Z",
      "short_description": "view는 query의 결과 집합을 하나의 가상의 테이블로 보여주는 방식이다. 즉, 실제의 테이블을 만드는 것이 아니라 가상의 테이블을 만드는 것이다. 핵심은 query에 우리의 이름을 부여하고 저장하는 query라는 것이다. 단, 결과로 나오는 테이블이 가상의 테이블",
      "thumnail": "",
      "url_slug": "SQL-재활-훈련-9일차-View와-Having",
      "tags": [
        "sql"
      ]
    }
  ]
}
```

3. Get post information in detail
- request
```sh
http://localhost:9596/v1/velog/chappi/post?url_slug=eBPF
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
curl http://localhost:9596/v1/velog/chappi/post/download?url_slug=eBPF
```

- response
```
post.zip
```

4. Download all post of user

- request
```sh
http://localhost:9596/v1/velog/chappi/posts/download
```

- response
```sh
chappi-velog-posts.zip
```

5. Download some post file

- request
```sh
curl -X POST localhost:9596/v1/velog/chappi/posts/download   -H "content-type: application/json" -d '[{"url_slug": "eBPF"}, {"url_slug": "SQL-재활-훈련-9일차-View와-Having"}]' --output result.zip
```

- response
```sh
zipfile.zip
```

6. get the all series of user

- request
```sh
curl localhost:9596/v1/velog/chappi/series
```

- response
```json
{
  "data": [
    {
      "id": "b9a6f3ba-31ec-48c1-b8a8-f86badbcfc25",
      "name": "eBPF",
      "count": 3,
      "thumbnail": "https://velog.velcdn.com/images/chappi/post/9eeb7980-69a2-43a8-a578-dad3b7e620f7/image.png",
      "updated_at": "2025-04-24T13:15:38.983Z"
    },
    {
      "id": "d86249cb-ab48-4127-b705-a5d158b02d5e",
      "name": "sql",
      "count": 9,
      "thumbnail": "",
      "updated_at": "2025-04-15T04:06:17.382Z"
    },
    {
      "id": "095574fd-fc24-490a-afdc-014076c27662",
      "name": "elasticsearch",
      "count": 9,
      "thumbnail": "",
      "updated_at": "2025-04-04T04:49:03.154Z"
    },
    ...
    {
      "id": "4c86cc6e-66fa-4dcc-b5e3-45b969ec6b0e",
      "name": "알고리즘",
      "count": 11,
      "thumbnail": "https://images.velog.io/images/chappi/post/a23b8edc-3e5a-4589-82f7-1e80727fb14f/멈춰!.jpg",
      "updated_at": "2021-04-28T05:26:25.247Z"
    }
  ]
}
```

7. read series in detail

- request
```sh
curl localhost:9596/v1/velog/chappi/series/b9a6f3ba-31ec-48c1-b8a8-f86badbcfc25
```

- response
```json

```