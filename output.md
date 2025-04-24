# Aggregation
aggregation(집계)는 elasticsearch의 꽃이다. search도 aggregation을 하기 위함과 다를바 없다. 

먼저 kibana sample data를 적재해보도록 하자. kibana homepage에 가서 `Try sample data` -> `Sample eCommerce orders`를 누르면 `kibana_sample_data_ecommerce` index가 적재된다.

elasticsearch의 aggregation은 search의 연장선이다. 즉 search 후에 그 결과를 바탕으로 aggregation하는 것이다. 다음은 search에 매칭된 document를 대상으로 지정한 field의 값을 모두 합친 값을 반환하는 aggregation요청이다.  
```json
GET kibana_sample_data_ecommerce/_search
{
  "size": 0,
  "query": {
    "term": {
      "currency": {
        "value": "EUR"
      }
    }
  },
  "aggs": {
    "my-sum-aggregation-name": {
      "sum": {
        "field": "taxless_total_price"
      }
    }
  }
}
```
요청을 보면 `search API`에 `aggs`가 추가된 것이 전부이다. `size`가 0인 것이 의문일텐데, `size`를 `0`으로 지정하면 search에 상위 매칭된 문서가 무엇인지 받아볼 수가 없다.

그러나 이와 상관없이 search조건에 매치되는 모든 document들은 aggregation작업에 사용되기 때문에, 문제가 없다. 또한, `size`를 0으로 지정하면 각 shards 대해서 `search`를 하지않고 aggregation만 할 뿐으로 성능에 이득도 있고 캐시에 도움도 더 많이 받을 수 있다.

`aggs`부분 아래에 실행할 aggregation의 이름을 적고 그 하위에 aggregation 종류를 기술한 후 필요한 값들을 넣어 요청한다. 

주의할 것은 aggregation은 search query에 매칭된 모든 document에 대해서 수행되므로, 과도한 양이 대상이 되지 않도록 해ㅐ야한다. 그런 경우 전체 클러스터의 성능을 급격히 저하시킬 수 있기 때문이다. 특히 kibana가 그런 측면에서 부담이 크다. 따라서 kibana는 팀 이외의 외부에 열어두지 않는 것이 좋다. 

응답은 다음과 같다.
```json
{
  "took" : 0,
  "timed_out" : false,
  "_shards" : {
    "total" : 1,
    "successful" : 1,
    "skipped" : 0,
    "failed" : 0
  },
  "hits" : {
    "total" : {
      "value" : 4675,
      "relation" : "eq"
    },
    "max_score" : null,
    "hits" : [ ]
  },
  "aggregations" : {
    "my-sum-aggregation-name" : {
      "value" : 350884.12890625
    }
  }
}
```
`value`부분이 `taxless_total_price`의 합이다. 

elasticsearch에서 지원하는 여러 aggregation을 크게보면 `metric`, `bucket`, `pipeline` aggregation으로 분류된다. 

## metric aggregation(집계)
metric aggregation은 document에 대한 산술적인 연산을 수행한다. 

### avg, max, min, sum aggregation
`avg`, `max`, `min`, `sum` aggregation은 search에 매칭된 document를 대상으로 지정한 field의 값을 가져온 뒤 각각 평균, 최대값, 최소값, 합을 계산하여 반환한다.
```json
GET kibana_sample_data_ecommerce/_search
{
  "size": 0,
  "query": {
    "term": {
      "currency": {
        "value": "EUR"
      }
    }
  },
  "aggs": {
    "my-sum-aggregation-name": {
      "avg": {
        "field": "taxless_total_price"
      }
    }
  }
}
```
`max`, `min`, `sum`도 `avg`가 있는 칸에 사용하면 된다.

### stats aggregation
`stats` aggregation은 지정한 field의 `avg`, `max`, `min`,`sum`, `개수`를 모두 계싼해서 반환한다.

```json
GET kibana_sample_data_ecommerce/_search
{
  "size": 0,
  "query": {
    "term": {
      "currency": {
        "value": "EUR"
      }
    }
  },
  "aggs": {
    "my-sum-aggregation-name": {
      "stats": {
        "field": "taxless_total_price"
      }
    }
  }
}
```
결과가 다음과 같이 나온다.

```json
{
...
  },
  "hits" : {
    "total" : {
      "value" : 4675,
      "relation" : "eq"
    },
    "max_score" : null,
    "hits" : [ ]
  },
  "aggregations" : {
    "my-sum-aggregation-name" : {
      "count" : 4675,
      "min" : 6.98828125,
      "max" : 2250.0,
      "avg" : 75.05542864304813,
      "sum" : 350884.12890625
    }
  }
}
```
이렇게 여러 숫자 값을 한꺼번에 반환하는 metric aggregation를 `multi-value numeric metric aggregation`이라고 부른다. 

### cardinality aggregation
`cardinality` aggregation은 지정한 field가 고유한 값의 개수를 계산해 반환한다. 
```json
GET kibana_sample_data_ecommerce/_search
{
  "size": 0,
  "query": {
    "term": {
      "currency": {
        "value": "EUR"
      }
    }
  },
  "aggs": {
    "my-cardinality-aggregation-name": {
      "cardinality": {
        "field": "customer_id",
        "precision_threshold": 3000
      }
    }
  }
}
```
위의 query는 `currency`의 `value`가 `EUR`이고 `cardinality`는 `customer_id`를 계산하라는 것이다. 즉, `customer_id`가 몇개가 있는 지 세어보라는 것과 같다. 

`precision_threshold` 옵션은 정확도를 조절하기 위해 사용한다. 이 값을 높이면 정확도가 올라가지만 그 만큼 메모리를 더 많이 사용한다. 다만 정확도를 올리기 위해 이 값을 무작정 많이 높일 필요는 없다. `precision_threshold`가 최종 `cardinality`보다 높다면 정확도가 충분히 높다. 반대로 `cardinality`값이 `precision_threshold`를 넘어서면 정확도가 떨어진다. 이를 감안해서 적당한 값을 지정해주는 것이 좋다.

default 로 3000이며 최대값은 40000이다. 

결과는 다음과 같다.
```json
{
  //...
  "hits" : {
    "total" : {
      "value" : 4675,
      "relation" : "eq"
    },
    "max_score" : null,
    "hits" : [ ]
  },
  "aggregations" : {
    "my-cardinality-aggregation-name" : {
      "value" : 46
    }
  }
}
```
`my-cardinality-aggregation-name.value` 결과는 `46`이므로 `customer_id`가 총 46개 있다는 것이다. 

## bucket aggregation
bucket aggregation은 document를 특정 기준으로 쪼개어 여러 부분 집합으로 나눈다. 이 부분 집합을 bucket이라고 한다. 또한 각 bucket에 포함된 문서를 대상으로 별도의 하위 aggregation(sub-aggregation)를 수행할 수 있다.

### range aggregation
`range` aggregation은 지정한 field값을 기준으로 document를 원하는 bucket구간으로 쪼갠다. bucket구간으로 나눌 기준이 될 field와 기준값을 지정해 요청한다.

먼저 sample data를 적재하기 위해서 browser main에 가서 `sample flight data`를 적재하도록 하자. `sample flight data`를 적재하였다면 `kibana_sample_data_flights`라는 index로 data가 적재된다.

다음으로 `range`를 써서 aggregation을 해보도록 하자.
```json
GET kibana_sample_data_flights/_search
{
  "size": 0,
  "query": {
    "match_all": {}
  },
  "aggs": {
    "distance-kilometers-range": {
      "range": {
        "field": "DistanceKilometers",
        "ranges": [
          {
            "to": 5000
          },
          {
            "from": 5000,
            "to": 10000
          },
          {
            "from": 10000
          }
        ]
      },
      "aggs": {
        "average-ticket-price": {
          "avg": {
            "field": "AvgTicketPrice"
          }
        }
      }
    }
  }
}
```
`range`는 `ranges`를 통해서 여러 bucket을 만들 수 있다. 이러한 bucket구간에 또 다른 aggregation을 바로 적용할 수 있는데, 그것이 바로 `average-ticket-price`로 `AvgTicketPrice`의 평균을 구하는 것이다. 재밌는 것은 `range`와 그 bucket에 대한 `aggs`가 같은 level에 있다는 것인데, 이는 `range` 후에 그 결과의 bucket이 바로 `aggs.average-ticket-price`로 간다고 생각하면 된다. 

결과는 다음과 같이 나온다.
```json
{
   , ///
  "aggregations" : {
    "distance-kilometers-range" : {
      "buckets" : [
        {
          "key" : "*-5000.0",
          "to" : 5000.0,
          "doc_count" : 4052,
          "average-ticket-price" : {
            "value" : 513.3930266305937
          }
        },
        {
          "key" : "5000.0-10000.0",
          "from" : 5000.0,
          "to" : 10000.0,
          "doc_count" : 6042,
          "average-ticket-price" : {
            "value" : 677.2621444606182
          }
        },
        {
          "key" : "10000.0-*",
          "from" : 10000.0,
          "doc_count" : 2965,
          "average-ticket-price" : {
            "value" : 685.3553124773563
          }
        }
      ]
    }
  }
}
```
`range`의 `ranges`에 나열한대로 3개의 bucket이 나온 것을 볼 수 있다. 이 bucekt에는 `average-ticker-price`가 있는데, 각 bucket이 `average-ticket-price`의 평균 aggregation연산이 적용된 결과를 볼 수 있다. 

만약, 위와 같이 bucket에 대한 하위 aggregation이 실행되지 않으면 `doc_count`까지만 세고 끝난다.

```json
GET kibana_sample_data_flights/_search
{
  "size": 0,
  "query": {
    "match_all": {}
  },
  "aggs": {
    "distance-kilometers-range": {
      "range": {
        "field": "DistanceKilometers",
        "ranges": [
          {
            "to": 5000
          },
          {
            "from": 5000,
            "to": 10000
          },
          {
            "from": 10000
          }
        ]
      }
    }
  }
}
```
이렇게 하위 aggregation에 대한 정의를 하지 않으면, 다음과 같이 `doc_count`만 세고 끝난다.

```json
{
    ,///
  "aggregations" : {
    "distance-kilometers-range" : {
      "buckets" : [
        {
          "key" : "*-5000.0",
          "to" : 5000.0,
          "doc_count" : 4052
        },
        {
          "key" : "5000.0-10000.0",
          "from" : 5000.0,
          "to" : 10000.0,
          "doc_count" : 6042
        },
        {
          "key" : "10000.0-*",
          "from" : 10000.0,
          "doc_count" : 2965
        }
      ]
    }
  }
}
```
이를 통해서 알 수 있는 것은 bucket aggregation의 핵심은 **하위 집계(aggregation)**에 있다는 것이다. document전체에 대해 하나의 aggregation을 수행ㅎ는 것이 아니라 이렇게 document를 여러 구간의 bucket으로 나눈 뒤 각 bucket에 대해서 하위 aggregation을 수행하도록 하는 것이다. 하위 aggregation에 또 bucket aggregation을 넣으면 다시 그 하위 aggregation을 지정하는 것도 가능하다. 그러나 하위 aggregation이 너무 깊어지면 성능에 심각한 문제가 생기니 적당히 지정해야한다.

### date_range aggregation
`date_range`는 `range`와 유사하나 `date` 타입 field를 대상으로 사용하며 `from`, `to`에 간단한 날짜 시간 계산식을 사용할 수 있다는 점에서 차이가 있다.
```json
GET kibana_sample_data_ecommerce/_search
{
  "size": 0,
  "query": {
    "term": {
      "currency": {
        "value": "EUR"
      }
    }
  },
  "aggs": {
    "date-range-aggs": {
      "date_range": {
        "field": "order_date",
        "ranges": [
          {
            "to": "now-10d/d"
          },
          {
            "from": "now-10d/d",
            "to": "now"
          },
          {
            "from": "now"
          }
        ]
      }
    }
  }
}
```
3개의 bucket을 만들되, `date`를 기준으로 만들 수 있는 것이다. 하나는 현재로 부터 10분전, 하나는 10분부터 지금까지, 하나는 지금부터를 말한다.

```json
{
    ...
  "aggregations" : {
    "date-range-aggs" : {
      "buckets" : [
        {
          "key" : "*-2023-12-02T00:00:00.000Z",
          "to" : 1.7014752E12,
          "to_as_string" : "2023-12-02T00:00:00.000Z",
          "doc_count" : 299
        },
        {
          "key" : "2023-12-02T00:00:00.000Z-2023-12-12T09:00:29.060Z",
          "from" : 1.7014752E12,
          "from_as_string" : "2023-12-02T00:00:00.000Z",
          "to" : 1.70237162906E12,
          "to_as_string" : "2023-12-12T09:00:29.060Z",
          "doc_count" : 1529
        },
        {
          "key" : "2023-12-12T09:00:29.060Z-*",
          "from" : 1.70237162906E12,
          "from_as_string" : "2023-12-12T09:00:29.060Z",
          "doc_count" : 2847
        }
      ]
    }
  }
}
```
각 shard에 aggregation요청이 분산되어 들어오면, elasticsearch는 그 내용을 shard요청 캐시에 올린다. 이후 동일한 aggregation요청이 같은 shard로 들어오면 이 shard요청 캐시의 데이터를 활용해 그대로 반환한다. 동일한 aggregation요청인지는 요청 본문이 동일한가로 구분한다. 그러나 `now`가 포함된 aggregation은 캐시되지 않는데, 이는 호출 시점에 따라 요청 내용이 달라지는 성격의 요청이기 때문이다.

새로운 데이터가 들어와서 index상태가 달라지면 shard요청 캐시는 무효화되기 때문에 고정된 index가 아니면 캐시 활용도가 떨어가진다.게다가 완전히 같은 aggregation을 여러 번 요청해야하는 상황도 많지 않기 때문에 이 점을 크게 신경 쓸 필요는 없다. 다만 `now`가 포함된 aggregation은 shard요청 캐시에 올라가지 않는다는 점을 인지하고 사용해야한다. 

### histogram aggregation
`histogram` aggregation는 지정한 field의 값을 기준으로 bucket을 나눈다는 점에서 `range` aggregation과 유사하다. 다른 점은 bucket 구분의 경계 기준값을 직접 지정하는 것이 아니라, bucket의 간격을 지정해서 경계를 나눈다는 점이다. 

```json
GET kibana_sample_data_flights/_search
{
  "size": 0,
  "query": {
    "match_all": {}
  },
  "aggs": {
    "my-histogram": {
      "histogram": {
        "field": "DistanceKilometers",
        "interval": 1000
      }
    }
  }
}
```
`interval`을 기준으로 bucket을 자동으로 나누어 생성해낸다. 0~1000(0<=x<1000), 1000~2000(1000<=x<2000) ... 이런 식이 되는 것이다. 응답은 다음과 같다.

```json
{
  
  "aggregations" : {
    "my-histogram" : {
      "buckets" : [
        {
          "key" : 0.0,
          "doc_count" : 1806
        },
        {
          "key" : 1000.0,
          "doc_count" : 1153
        },
        {
          "key" : 2000.0,
          "doc_count" : 530
        },
        {
          "key" : 3000.0,
          "doc_count" : 241
        },
        ...
      ]
    }
  }
}
```

`interval`만을 입력하면 0을 시작으로 histogram 계급을 나눈다. 즉 `[0, 1000), [1000, 2000)` 이런식으로 생성되는 것이다. 시작 위치를 바꾸고 싶다면 `offset`을 사용하면 된다.
```json
GET kibana_sample_data_flights/_search
{
  "size": 0,
  "query": {
    "match_all": {}
  },
  "aggs": {
    "my-histogram": {
      "histogram": {
        "field": "DistanceKilometers",
        "interval": 1000,
        "offset": 50
      }
    }
  }
}
```
`offset`이 50이기 때문에, 50 - 1000부터 시작한다. 즉, `[-950, 50), [50, 1050), [1050, 2050)`으로 bucket이 만들어지는 것이다. 

```json
#! Elasticsearch built-in security features are not enabled. Without authentication, your cluster could be accessible to anyone. See https://www.elastic.co/guide/en/elasticsearch/reference/7.17/security-minimal-setup.html to enable security.
{
    ...
  "aggregations" : {
    "my-histogram" : {
      "buckets" : [
        {
          "key" : -950.0,
          "doc_count" : 643
        },
        {
          "key" : 50.0,
          "doc_count" : 1231
        },
        {
          "key" : 1050.0,
          "doc_count" : 1111
        },
        {
          "key" : 2050.0,
          "doc_count" : 522
        },
        ...
      ]
    }
  }
}
```
참고로 `DistanceKilometers` field는 음수값이 존재하지않지만 `[-950, 50)`구간을 통해서 0~49까지의 값을 측정할 수 있는 것이다.

이 밖에도 `min_doc_count`를 지정해서 bucket내 document수가 일정 이하인 bucket은 결과에서 제외할 수 있다. 

### date_histogram
`date_histogram` aggregation은 `histogram` aggregation과 유사하지만 대상으로 `date` type을 사용한다는 점에서 다르다. 또한, `interval`대신에 `calendar_interval`이나 `fixed_interval`을 사용한다. 
```json
GET kibana_sample_data_ecommerce/_search
{
  "size": 0,
  "query": {
    "match_all": {}
  },
  "aggs": {
    "my-date-histogram": {
      "date_histogram": {
        "field": "order_date",
        "calendar_interval": "day"
      }
    }
  }
}
```
`date_histogram`에서 `calendar_interval`을 `day`로 사용하면 document를 bucket단위로 쪼개어낸다. 

```json
#! Elasticsearch built-in security features are not enabled. Without authentication, your cluster could be accessible to anyone. See https://www.elastic.co/guide/en/elasticsearch/reference/7.17/security-minimal-setup.html to enable security.
{
  //
  "aggregations" : {
    "my-date-histogram" : {
      "buckets" : [
        {
          "key_as_string" : "2023-11-30T00:00:00.000Z",
          "key" : 1701302400000,
          "doc_count" : 146
        },
        {
          "key_as_string" : "2023-12-01T00:00:00.000Z",
          "key" : 1701388800000,
          "doc_count" : 153
        },
        {
          "key_as_string" : "2023-12-02T00:00:00.000Z",
          "key" : 1701475200000,
          "doc_count" : 143
        },
        ...
      ]
    }
  }
}
```

`calendar_interval`에는 다음과 같은 값들을 지정할 수 있다.
1. `minute` 또는 1m: 분 단위
2. `hour` 또는 1h: 시간 단위
3. `day` 또는 1d: 일 단위
4. `month` 또는 1M: 월 단위
5. `quarter` 또는 1q: 분기 단위
6. `year` 또는 1y: 연 단위

`calendar_interval`은 구체적인 시간을 요청할 수 없다. 즉, `1m`, `1h`은 가능하지만 `12m`, `13h`와 같은 시간은 불가능하다. 

이 경우에는 `fixed_interval `를 사용해야한다. `fixed_interval`는 `ms`, `s`, `m`, `h`, `d` 단위로 사용할 수 있기 때문이다. 가령 3시간 단위로 bucket을 만들고 싶다면 `fixed_interval: 3h`으로 만들면 된다.

추가적으로 `date_histogram` aggregation도 `offset`과 `min_doc_count`를 설정할 수 있다.

### terms aggregation
`terms` aggregation은 지정한 field에 대해 가장 빈도수가 높은 `term`순서대로 bucket을 생성한다. bucket을 최대 몇 개까지 생성할 것인지를 `size`로 지정한다.

test를 위해서 sample data로 `Sample web logs`를 다운받도록 하자.
```json
GET kibana_sample_data_logs/_search
{
  "size": 0,
  "query": {
    "match_all": {}
  },
  "aggs": {
    "my-terms-aggs": {
      "terms": {
        "field": "host.keyword",
        "size": 10
      }
    }
  }
}
```
위의 경우는 `host.keyword` field에 대해서 `term`의 개수를 새고 많은 수부터 bucket을 나열한다. 

```json
{
    //...
  },
  "aggregations" : {
    "my-terms-aggs" : {
      "doc_count_error_upper_bound" : 0,
      "sum_other_doc_count" : 0,
      "buckets" : [
        {
          "key" : "artifacts.elastic.co",
          "doc_count" : 6488
        },
        {
          "key" : "www.elastic.co",
          "doc_count" : 4779
        },
        {
          "key" : "cdn.elastic-elastic-elastic.org",
          "doc_count" : 2255
        },
        {
          "key" : "elastic-elastic-elastic.org",
          "doc_count" : 552
        }
      ]
    }
  }
}
```
`terms` aggregation은 각 shard에서 `size`개수만큼 `term`을 뽑아 빈도수를 센다. 각 shard에서 수행된 계산을 한 곳으로 모아 합산한 후 `size`개수만큼 bucket을 뽑는다. 그러므로 `size`개수와 각 document의 분포에 따라 그 결과가 정확하지 않을 수 있다. 각 bucket의 `doc_count`는 물론 하위 aggregation 결과도 정확한 수치가 아닐 수 있다. 특히, 해당 field의 고유한 `term`개수가 `size`보다 많다면 상위에 뽑혀야 할 `term`이 최종 결과에 포함되지 않을 수 있다.

응답의 `doc_count_error_upper_bound` field는 `doc_count`의 오차 상한선을 나타낸다. 이 값이 크다면 `size`를 높이는 것을 고려할 수 있다. 물론 `size`를 높이면 정확도는 올라가지만 그만큼 성능이 하락한다. 

`sum_other_doc_count` field는 최종적으로 bucket에 포함되지 않은 document 수를 나타낸다. 상위 `term`에 들지 못한 document 개수의 총합이다. 

만약 모든 `term`에 대해서 pagenation으로 전부 순회하며 aggregation을 하려고 한다면 `size`를 무작정 계속 높이는 것보다는 `composite` aggregation을 사용하는 것이 좋다. `terms` aggregation는 기본적으로 상위 `term`을 뽑아서 aggregation을 수행하도록 설계되었다. 

### composite aggregation
`composite` aggregastion은 `sources`로 지정된 하위 aggregation의 bucket을 pagenation을 이용해서 효율적으로 순회하는 aggregation이다. `sources`에 하위 aggregation를 여러 개 지정한 뒤 조합된 bucket을 생성할 수 있다.

```json
GET kibana_sample_data_logs/_search
{
  "size": 0,
  "query": {
    "match_all": {}
  },
  "aggs": {
    "composite-aggs": {
      "composite": {
        "size": 100, 
        "sources": [
          {
            "terms-aggs": {
              "terms": {
                "field": "host.keyword"
              }
            }
          },
          {
            "date-histogram-aggs": {
              "date_histogram": {
                "field": "@timestamp",
                "calendar_interval": "day"
              }
            }
          }
        ]
      }
    }
  }
}
```
`composite`아래의 `size`는 한 번에 몇 개의 bucket을 반환할 것인가를 지정한다. `sources`에는 bucket을 조합하여 순회할 하위 aggregation을 지정한다. 여기에는 모든 종류의 aggregation을 하위 aggregation으로 지정할 수는 없다. `terms` aggregation, `histogram` aggregation, `date_histogram` aggregation 등 일부 aggregation만 지정할 수 있다. 

잘보면 `terms` aggregation에는 `size`가 없다. 이는 `composite`자체가 bucket을 순차적으로 방문하는 목적의 aggregation이기 때문에 `terms`에 `size`개념이 필요하지 않다. 즉 모든 `term`을 뽑아낸다는 것이다.

위의 예시는 100개의 bucket을 만들되 `host.keyword`의 `term`중 상위로 나오는 `term`의 100개로 뽑아내고, 이를 `@timestamp`기준으로 `day`마다 aggregation하는 것이다.

```json
{
  //..
  "aggregations" : {
    "composite-aggs" : {
      "after_key" : {
        "terms-aggs" : "cdn.elastic-elastic-elastic.org",
        "date-histogram-aggs" : 1704844800000
      },
      "buckets" : [
        {
          "key" : {
            "terms-aggs" : "artifacts.elastic.co",
            "date-histogram-aggs" : 1701561600000
          },
          "doc_count" : 124
        },
        // ...
        {
          "key" : {
            "terms-aggs" : "cdn.elastic-elastic-elastic.org",
            "date-histogram-aggs" : 1704758400000
          },
          "doc_count" : 37
        },
        {
          "key" : {
            "terms-aggs" : "cdn.elastic-elastic-elastic.org",
            "date-histogram-aggs" : 1704844800000
          },
          "doc_count" : 24
        }
      ]
    }
  }
}
```
`artifacts.elastic.co` term이 가장 많이 나온 것이고, 그 다음이 `cdn.elastic-elastic-elastic.org`이다. 다음으로 해당 `term`을 가진 document에 대해서 날짜기준으로 aggregation한 것이다. 

`after_key`부분에서 확인할 수 있는 조합이 바로 pagenation을 위해서 필요한 가장 마지막 bucket의 key다. 이 `after_key`를 가져와서 다음과 같이 요청하면 작업 결과의 다음 페이지를 조회할 수 있다. 처음 요청과 동일한 내용이지만 `composite`아래에 `after`부분이 추가되었다.
```json
GET kibana_sample_data_logs/_search
{
  "size": 0,
  "query": {
    "match_all": {}
  },
  "aggs": {
    "composite-aggs": {
      "composite": {
        "size": 100, 
        "sources": [
          {
            "terms-aggs": {
              "terms": {
                "field": "host.keyword"
              }
            }
          },
          {
            "date-histogram-aggs": {
              "date_histogram": {
                "field": "@timestamp",
                "calendar_interval": "day"
              }
            }
          }
        ],
        "after": {
          "terms-aggs" : "cdn.elastic-elastic-elastic.org",
          "date-histogram-aggs" : 1704844800000
        }
      }
    }
  }
}
```
`after`부분에 응답의 `after_key`을 그대로 써주면 다음 pagenation이 가능한 것이다. 