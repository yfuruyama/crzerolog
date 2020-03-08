IMAGE=gcr.io/${PROJECT_ID}/crzerolog-cloudrun-http

build:
	docker build -t ${IMAGE} -f Dockerfile ../../

push:
	docker push ${IMAGE}

deploy:
	gcloud run deploy --image=${IMAGE} --platform=managed --region=us-central1 --project=${PROJECT_ID} --allow-unauthenticated crzerolog-cloudrun-http
