curl --get --data-urlencode 'text=Un viaje puede tener diversas motivaciones' http://localhost:8080/detect

curl --get --data-urlencode 'text=sentence to detect that I hope' http://localhost:8080/detect | pjson

curl --get --data-urlencode 'text=para guardar los momentos inolvidables del viaje' http://localhost:8080/detect

curl --get --data-urlencode 'text=Al llegar al lugar de destino es recomendable que te manejes llevando una mochila donde llevarás alimento, la cámara y tus documentos' http://localhost:8080/detect

curl --get --data-urlencode 'text=Panini-Bildchen sind für deutsche Kinder das, was Baseballkarten für kleine' http://localhost:8080/detect

curl --get --data-urlencode 'text=In Deutschland gibt es eine Institution, die zwar sehr wichtig ist, aber bei Nicht-Deutschen nicht so bekannt ist: die Stiftung Warentes' http://localhost:8080/detect

curl --get --data-urlencode 'text=sind bekannt für ihren Ordnungs- und Sauberkeitssinn' http://localhost:8080/detect
