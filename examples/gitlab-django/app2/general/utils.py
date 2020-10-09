import random
import os
import io
import boto3
"""
This function return random ASCII-string
"""
def generate_random_string():
    out_str = ''
    for i in range(0, 16):
        a = random.randint(65, 90)
        out_str += chr(a)
    return(out_str)

def upload_file(file_name, bucket, object_name=None):
    """Upload a file to an S3 bucket

    :param file_name: File to upload
    :param bucket: Bucket to upload to
    :param object_name: S3 object name. If not specified then file_name is used
    :return: True if file was uploaded, else False
    """
    s3_client = boto3.client(
        's3',
        endpoint_url = "https://s3.selcdn.ru",
        aws_access_key_id = os.environ.get("SPACE_ACCESS"),
        aws_secret_access_key = os.environ.get("SPACE_SECRET"),
        region_name = "ru-1a"
    )

    if object_name is None:
        object_name = file_name.split("/")[-1]

    with open(file_name, "rb") as f:
        s3_client.upload_fileobj(f, bucket, object_name)