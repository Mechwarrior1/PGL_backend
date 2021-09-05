# GORecycle Frontend

This is a course project, a functioning prototype website

## The Idea
A marketplace for people to post functional items to be given away, an organisation can take for further distribution or an individual can then pick them up. Helps to reduce waste and improve opportunities for reuse of items/materials


## Heroku and Frontend
hosted on heroku: http://gorecycle.herokuapp.com/
Frontend: https://github.com/Mechwarrior1/PGL_frontend


## API
| API | Params | Query | Payload | Return | Comments |
| --- | --- | --- | --- | --- | --- |
| /api/v0/health | None | None | None | StatusOK| see if server is online |
| /api/v0/ready | None | None | None | StatusOK/Status Unavailable | see if server is ready for traffic |
| /api/v0/db/info PUT | None | None | Json payload with Key and DataInfo | StatusCreated/StatusBadRequest | used to edit information on database |
| /api/v0/db/info GET | None | id & db | None | StatusCreated/StatusBadRequest with json payload | get information on database |
| /api/v0/db/completed/ PUT | None | id | Json payload with Key and Username | StatusCreated/StatusBadRequest | Change ItemListing entry on database to completed |
| /api/v0/db/signup POST | None | None | Json payload with with user info | StatusCreated/StatusBadRequest | Log new user into database |
| /api/v0/db/info POST | None | None | Json payload with with entry info | StatusCreated/StatusBadRequest | Log new entry into database |


credits:
word2vec codes: https://github.com/danieldk/go2vec/blob/8029f60947ae/go2vec.go
word2vec slim bin file: https://github.com/eyaler/word2vec-slim
