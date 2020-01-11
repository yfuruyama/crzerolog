IMAGE=gcr.io/${PROJECT_ID}/crzerolog-example

build:
	docker build -t ${IMAGE} .

push: build
	docker push ${IMAGE}

deploy: push
	gcloud run deploy --image=${IMAGE} --platform=managed --region=us-central1 --project=${PROJECT_ID} --allow-unauthenticated crzerolog-example
