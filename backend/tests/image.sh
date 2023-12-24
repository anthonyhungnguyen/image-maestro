# Extract EXIF
curl -X POST -H "Authorization: Bearer 123" --data '{"url": "https://media.npr.org/assets/img/2022/05/25/gettyimages-917452888-edit_custom-c656c35e4e40bf22799195af846379af6538810c-s1100-c50.jpg","exif":true}' http://localhost:8080/api/image | jq

# Annotations
curl -X POST -H "Authorization: Bearer 123" http://localhost:8080/api/image/123/annotate | jq