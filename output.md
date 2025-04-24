# template
template를 사전에 정의하면 index생성 시에 template에 정의한대로 index가 생성된다. 즉, mapping이나 type설정들에 대해서 template를 먼저 만들어 놓고, 생성되는 index가 해당 template의 정의를 따르도록 할 수 있다는 것이다. 이를 통해 업무 효율을 향상시키고 반복 작업이 야기할 수 있는 사람의 실수를 줄여 준다. 

## index template
index template는 index의 이름에 pattern을 적용하여 `wildcard(*)`와 같은 pattern 매칭을 시킬 수 있다. 생성된 index는 index template의 pattern에 일치하면 template의 제약을 받게 된다. 가령 index template의 `pattern`이 `park*`이고, 아무것도 없이 생성된 index가 `park-12`라면, 적용을 받는 것이다.

그러나, index template가 두 개가 있어서 하나의 index가 두 개의 index template의 영향을 받으면 어떻게해야할까?? index template가 `park*`과 `*-12`라면 index `park-12`는 두 index template의 영향을 받게된다. 이러한 일을 해결하기 위하여 template에서는 `priority`를 제공하는데, `priority`가 높을 수록 우선순위가 높아 먼저 적용받는다.

template는 다음과 같이 만들 수 있다.
```json
PUT _index_template/my_template
{
  "index_patterns": [
    "pattern_test_index-*",
    "another_pattern-*"
  ],
  "priority": 1,
  "template": {
    "settings": {
      "number_of_shards": 2,
      "number_of_replicas": 2
    },
    "mappings": {
      "properties": {
        "myTextField": {
          "type": "text"
        }
      }
    }
  }
}
```
이른은 `my_template`이고 index pattern은 `pattern_test_index-*`와 `another_pattern-*`이다. template를 통해서 shard수도 조정가능하고 `mappings`를 설정해놓아서 특정 properties를 설정할 수 있다. `myTextField`는 `type`이 `text`이므로, 해당 template를 따르는 index는 `myTextField`가 `text`이며 shard수는 2개, replica는 2개가 된다.

이제 index를 만들어보도록 하자.
```json
PUT pattern_test_index-1
```

`pattern_test_index-1`은 `pattern_test_index-*` pattern을 지키기 때문에 문제없이 template의 영향을 받는다. 따라서 위의 template설정이 유효하게 적용된다.
```json
GET pattern_test_index-1

...
{
  "pattern_test_index-1" : {
    "aliases" : { },
    "mappings" : {
      "properties" : {
        "myTextField" : {
          "type" : "text"
        }
      }
    },
    "settings" : {
      "index" : {
        "routing" : {
          "allocation" : {
            "include" : {
              "_tier_preference" : "data_content"
            }
          }
        },
        "number_of_shards" : "2",
        "provided_name" : "pattern_test_index-1",
        "creation_date" : "1702015818382",
        "number_of_replicas" : "2",
        "uuid" : "3EDW8Sw1RACo1z7zd7sRqA",
        "version" : {
          "created" : "7171499"
        }
      }
    }
  }
}
```
properties수로 `myTextField`를 가지고 있고 `type`도 `text`이다. `shards`수는 2개이며 `replicas`수도 2개이다. 

## component template
template를 만들다보면 여러 부분들이 template간에 겹치는 것을 확인할 수 있다. 가령 위의 `shards`, `replicas`수와 `mappings` 정보들이 그렇다. 그래서 이러한 부분들을 여러 개의 component들로 나눈다음, index-template를 만들 때 조합하는 것이다.

```json
PUT _component_template/timestamp_mappings
{
  "template": {
    "mappings": {
      "properties": {
        "timestamp": {
          "type": "date"
        }
      }
    }
  }
}

PUT _component_template/my_shard_settings
{
  "template": {
    "settings": {
      "number_of_shards": 2,
      "number_of_replicas": 2
    }
  }
}
```
두 개의 component template를 만들었다. 하나는 `mapping`정보에 대한 정의를 담고 있고, 하나는 `setting`정보에 대한 정의를 담고 있다.

이제 index template를 만들 때 다음의 두 component를 조합해 하나의 index template를 만들 수 있다.
```json
PUT _index_template/my_template2
{
  "index_patterns": ["timestamp_index-*"],
  "composed_of": ["timestamp_mappings", "my_shard_settings"]
}
```

이제 해당 template pattern을 만족하는 index template를 만들어보도록 하자.
```sh
PUT timestamp_index-001
GET timestamp_index-001
```

다음의 결과를 볼 수 있다.
```json
{
  "timestamp_index-001" : {
    "aliases" : { },
    "mappings" : {
      "properties" : {
        "timestamp" : {
          "type" : "date"
        }
      }
    },
    "settings" : {
      "index" : {
        "routing" : {
          "allocation" : {
            "include" : {
              "_tier_preference" : "data_content"
            }
          }
        },
        "number_of_shards" : "2",
        "provided_name" : "timestamp_index-001",
        "creation_date" : "1702017043845",
        "number_of_replicas" : "2",
        "uuid" : "cDGXouBfS9SCc3-urcLnxw",
        "version" : {
          "created" : "7171499"
        }
      }
    }
  }
}
```
`timestamp_mappings`와 `my_shard_settings` component template 내용이 모두 잘 적용된 것을 볼 수 있다. 이렇게 index template를 만들고 관리하면 훨씬 더 깔끔하게 index template를 관리할 수 있다.

## legacy template
index template와 component template API는 elasticsearch 7.8버전부터 추가된 기능이다. 이전에 사용하던 template API는 legacy template로 취급된다. 이전 버전의 template 기능은 `_index_template`대신에 `_template`을 사용하며, component template로 조합할 수 없다는 점을 제외하면 거의 동일하다.

따라서, 정의하는 방법 역시도 동일하다. 다만, legacy template는 적용 우선순위가 새 index template에 비해서 낮으므로 조심해야한다. 즉, 매칭되는 index template가 없을 때만 legacy_template를 적용한다는 것이다.

## dynamic template
dynamic template는 매칭되는 index의 properties의 이름을 보고 mapping을 dynamic하게 지정해주도록 한다. 가령, `dynamic_mapping-1`이라는 index에서 properties가 접두사로 `text`를 가지면 `text` type으로 `keyword`를 가지면 `keyword`로 type을 지정하도록 하는 것이다. 가령, `book_text`라는 property를 가지면 `text` type을 가지도록 하고 `title_keyword`라는 property가 오면 `keyword` type을 가지도록 하는 것이다. 

```json
PUT _index_template/dynamic_mapping_template
{
  "index_patterns": ["dynamic_mapping*"],
  "priority": 1,
  "template": {
    "settings": {
      "number_of_shards":  2,
      "numberof_replicas": 2
    }
  },
  "mappings": {
    "dynamic_templates": [
      {
        "my_text": {
          "match_mapping_type": "string",
          "match": "*_text",
          "mapping": {
            "type": "text"
          }
        }
      },
      {
        "my_keyword": {
          "match_mapping_type": "string",
          "match": "*_keyword",
          "mapping": {
            "type": "keyword"
          }
        }
      }
    ]
  }
}
```
위의 `dynamic_mapping_template`는 `dynamic_mapping*` index matching을 통과하는 index들에 한해서, index의 property가 `name_text`면 `name_text`에 `text` field를 부여하고, `name_keyword`면 `name_keyword` property에 `keyword` type을 부여한다. 

이렇게 동적인 mapping을 지원할 수 있는 것은 `dynamic_templates`의 `match`부분으로 `match`의 조건을 만족하는 property는 `mapping`에 해당하는 것이다. 이러한 `match`조건들은 다음과 같다.

1. `match_mapping_type`: 새로 들어오는 data type을 `JSON` parser를 이용해 확인한다. `JSON` parser는 `long`과 `interger`등을 인지할 수 없기 때문에 data type지정에 제한이 된다. `boolean`, `double`, `long`, `string`, `object`, `date`

2. `match/unmatch`: property이름이 지정된 패턴과 일치할 때, 또는 일치하지 않을 때 적용한다.

3. `path_match / path_unmatch`: `match/unmatch`와 동일하게 동작하지만 field이름으로 마침표를 사용한 전체 경로를 이용한다. `my_object.name.text*`
