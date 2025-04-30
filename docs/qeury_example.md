1. Get posts with limit

- request
```sh
http://localhost:9596/v1/velog/posts?name=chappi&count=100&post_id=36181f27-fcb7-4164-89c4-5d6db6c1b2ee
```

- response
```json
{
  "posts": [
    {
      "id": "aeaf24f4-6a15-4fef-9d1b-ac98bb53d82b",
      "title": "Java 재활 훈련 9일차 - Object class, record, System properties, StringBuilder, StringTokenizer, Wrapper class, reflection, Annotation",
      "created_at": "2025-02-28T08:38:04.583Z",
      "updated_at": "2025-04-11T04:36:13.286Z",
      "short_description": "클래스를 선언할 때 extends 키워드로 다른 클래스를 상속하지 않으면 암묵적으로 java.lang.Object 클래스를 상속하게 된다. 따라서, 모든 클래스는 Object의 자손 클래스이다. 모든 객체들은 Object가 가진 메서드를 사용할 수 있다.equals()",
      "thumnail": "",
      "url_slug": "Java-재활-훈련-9일차-Object-class-record-System-properties-StringBuilder-StringTokenizer-Wrapper-class-reflection-Annotation",
      "tags": [
        "Java"
      ]
    },
    {
      "id": "1cfef770-b342-46ab-8865-84f730cc14c0",
      "title": "Java 재활 훈련 8일차 - Exception",
      "created_at": "2025-02-28T08:36:12.266Z",
      "updated_at": "2025-04-14T07:15:52.643Z",
      "short_description": "java에서는 예와(exception)라고 부르는 오류가 있다. exception은 잘못된 문법, 코딩으로 인해 발생한 오류를 말한다. exception이 발생하면 프로그램은 곧바로 종료된다는 점에서 에러와 동일하지만, exception의 처리를 통해 실행 상태를 유지",
      "thumnail": "",
      "url_slug": "Java-재활-훈련-8일차-Exception",
      "tags": [
        "Java"
      ]
    },
    ...
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
http://localhost:9596/v1/velog/post/download?name=chappi&url_slog=eBPF
```

- response
```
post.zip
```