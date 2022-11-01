import boto3
from os import environ

endpoint_url = 'http://127.0.0.1:9900'
access_key = 'minioadmin'
secret_key = 'minioadmin'
bucket_name = 'werf-guide-app'
filename = 'filename'

if environ.get('MINIO_ACCESS_KEY') is not None:
    access_key = environ.get('MINIO_ACCESS_KEY')

if environ.get('MINIO_SECRET_KEY') is not None:
    secret_key = environ.get('MINIO_SECRET_KEY')

if environ.get('BUCKET_NAME') is not None:
    bucket_name = environ.get('BUCKET_NAME')

if environ.get('MINIO_ENDPOINT') is not None:
    endpoint_url = environ.get('MINIO_ENDPOINT')

s3_client = boto3.client('s3',
                         aws_access_key_id=access_key,
                         aws_secret_access_key=secret_key,
                         endpoint_url=endpoint_url,
                         config=boto3.session.Config(signature_version='s3v4')
                         )


def create_bucket_if_not_exists(bucket_name):
    for b in s3_client.list_buckets()["Buckets"]:
        if b["Name"] == bucket_name:
            return
    s3_client.create_bucket(Bucket=bucket_name)


def upload_to_s3(content):
    create_bucket_if_not_exists(bucket_name)
    s3_client.put_object(Bucket=bucket_name, Key=filename, Body=content)


def download_from_s3():
    try:
        response = s3_client.get_object(Bucket=bucket_name, Key=filename)
    except:
        return "You haven't uploaded anything yet.\n"
    return response['Body'].read()
