# 📦 Velog Story Shift

> **Download your posts from Velog (or other markdown-based platforms) and migrate them to another markdown blog with ease.**

---

## ✨ About

**Velog Story Shift** is an application that helps you download posts from markdown-based blogging platforms like Velog and migrate them to another markdown-compatible blog with ease.

Even when switching blog platforms, you can preserve your content in clean, readable markdown format.

---

## 🚀 Objective
- Fetch posts using the Velog API  
- Save posts as Markdown files  
- Migration-ready structure for compatibility with other platforms  
- Customizable output directory and file naming  
- Upcoming support: Tistory, Hugo, Gatsby, Jekyll, and more  

---

## 🔧 Usage
1. Start a server by script
```bash
./script/start.sh -c web
```

or you can use cli version

```bash
./script/start.sh -c cli
```

2. Start a server by go
```bash
go run ./...
```

3. run test
```bash
go test ./... -v
```

4. test query script
```bssh
 ./script/query_test.sh
```

5. dockerfile

- build and run
```bssh
./script/docker_script.sh -t 1.0.1 -f ./Dockerfile -R "-p 9596:9596 -e STORY_SHIFT_CONFIG_FILE=./config/test_config.yaml"
```

- build
```bash
./script/docker_script.sh -t 1.0.1 -f ./Dockerfile
```

- run
```bssh
./script/docker_script.sh -t 1.0.1 -f ./Dockerfile -c run -R "-p 9596:9596 -e STORY_SHIFT_CONFIG_FILE=./config/test_config.yaml"
```

---

# 📘 Velog REST API

Velog 데이터를 REST 방식으로 제공하는 API입니다.

| **Operation**                                   | **Request URL**                                                                                       | **Request Type** | **Request Payload**                                                                                                      |
|-------------------------------------------------|------------------------------------------------------------------------------------------------------|------------------|--------------------------------------------------------------------------------------------------------------------------|
| **1. Get posts with limit**                     | `http://localhost:9596/v1/velog/{username}/posts?count=10`                                                 | GET              | No payload                                                                                                               |
| **2. Get posts with limit and post_id**         | `http://localhost:9596/v1/velog/{username}/posts?count=2&post_id={post_id}`    | GET              | No payload                                                                                                               |
| **3. Get post information in detail**           | `http://localhost:9596/v1/velog/{username}/post?url_slug={url_slug}`                                            | GET              | No payload                                                                                                               |
| **4. Download post**                            | `curl http://localhost:9596/v1/velog/{username}/post/download?url_slug={url_slug}`                              | GET              | No payload                                                                                                               |
| **5. Download all posts of user**               | `http://localhost:9596/v1/velog/{username}/posts/download`                                                | GET              | No payload                                                                                                               |
| **6. Download some post files**                 | `curl -X POST localhost:9596/v1/velog/{username}/posts/download -H "content-type: application/json" -d '[{"url_slug": "eBPF"}, {"url_slug": "SQL-재활-훈련-9일차-View와-Having"}]' --output result.zip` | POST             | Payload: `[{"url_slug": "eBPF"}, {"url_slug": "SQL-재활-훈련-9일차-View와-Having"}]`                                      |
| **7. Get all series of user**                   | `curl localhost:9596/v1/velog/{username}/series`                                                          | GET              | No payload                                                                                                               |
| **8. Get all posts in selected series**         | `curl http://localhost:9596/v1/velog/{username}/series/{series_slug}`                                              | GET              | No payload                                                                                                               |
| **9. Get all posts in series as zip file**     | `curl http://localhost:9596/v1/velog/{username}/series/{series_slug}/download`                                              | GET              | No payload                                                                                                               |



## ✅ Example
---


## 🛠 Configuration
You can set config file path by using 'STORY_SHIFT_CONFIG_FILE' env
```
export STORY_SHIFT_CONFIG_FILE=config.yaml
```

We now support yaml, env configuration data, please check the config directory


---

## 🤝 Contribution

Feel free to open issues, suggest features, or submit pull requests!

---