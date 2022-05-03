GOCMD=go
GOBUILD=${GOCMD} build -mod=mod
GOCLEAN=${GOCMD} clean

build: api

.PHONY: \
    api social feed article comment user search chat async

clean:
	${GOCLEAN}

api:
	${GOBUILD} -o /home/work/run/zzlove_api zzlove/cmd/api

social:
	${GOBUILD} -o /home/work/run/zzlove_social zzlove/cmd/social

feed:
	${GOBUILD} -o /home/work/run/zzlove_feed zzlove/cmd/feed

article:
	${GOBUILD} -o /home/work/run/zzlove_article zzlove/cmd/article

comment:
	${GOBUILD} -o /home/work/run/zzlove_comment zzlove/cmd/comment

user:
	${GOBUILD} -o /home/work/run/zzlove_user zzlove/cmd/user

search:
	${GOBUILD} -o /home/work/run/zzlove_search zzlove/cmd/search

chat:
	${GOBUILD} -o /home/work/run/zzlove_chat zzlove/cmd/chat

async:
	${GOBUILD} -o /home/work/run/zzlove_async zzlove/cmd/async
