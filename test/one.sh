#!/bin/bash

# 开篇立论
curlie -d '{"topic": "爱与被爱，哪个更幸福？", "role":"1", "question":"辩题是爱与被爱，哪个更幸福？，首先是立论环节，请正方开篇立论。"}' \
    POST \
    http://47.96.12.150:8080/memory/create

# 自由辩论
curlie -d '{"dialog": { "question": "被爱更幸福啦，你懂不了的啦", "answer": "你谁？" }, "last": false }' \
    POST \
    http://47.96.12.150:8080/memory/2/update
