# ocean

Pure memory NewSQL based on raft protocol  (Born for High Frequency Trading)  

> Powered by 成都珠码开天信息科技有限公司


提供类似elasticsearch的DSL查询

```sql 
{
    "param":[
        {
            "filter_type":"gt",
            "key":"age",
            "value":30,
            "params":null
        },
        {
            "filter_type":"like",
            "key":"name",
            "value":"wamg",
            "params":null
        },
        {
            "filter_type":"must",
            "key":"",
            "value":null,
            "params":[
                {
                    "filter_type":"gt",
                    "key":"age",
                    "value":20,
                    "params":null
                },
                {
                    "filter_type":"gt",
                    "key":"money",
                    "value":60,
                    "params":null
                },
                {
                    "filter_type":"must",
                    "key":"",
                    "value":null,
                    "params":[
                        {
                            "filter_type":"gt",
                            "key":"money",
                            "value":180,
                            "params":null
                        }
                    ]
                }
            ]
        }
    ],
    "filter_type":"must",
    "from":0,
    "size":10,
    "sort":[
        {
            "key":"age",
            "sort_type":"desc"
        }
    ]
}
```