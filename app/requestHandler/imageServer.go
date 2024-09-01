package requestHandler

import (
	"golang.org/x/net/context"
	"io/ioutil"
	"lifesgood/app/config"
	"log"
	"net/http"
	"strings"

	"cloud.google.com/go/storage"
)

// createFile creates a file in Google Cloud Storage.
//func (d *demo) createFile(fileName string) {
//	fmt.Fprintf(d.w, "Creating file /%v/%v\n", d.bucketName, fileName)
//
//	wc := d.bucket.Object(fileName).NewWriter(d.ctx)
//	wc.ContentType = "text/plain"
//	wc.Metadata = map[string]string{
//		"x-goog-meta-foo": "foo",
//		"x-goog-meta-bar": "bar",
//	}
//	d.cleanUp = append(d.cleanUp, fileName)
//
//	if _, err := wc.Write([]byte("abcde\n")); err != nil {
//		d.errorf("createFile: unable to write data to bucket %q, file %q: %v", d.bucketName, fileName, err)
//		return
//	}
//	if _, err := wc.Write([]byte(strings.Repeat("f", 1024*4) + "\n")); err != nil {
//		d.errorf("createFile: unable to write data to bucket %q, file %q: %v", d.bucketName, fileName, err)
//		return
//	}
//	if err := wc.Close(); err != nil {
//		d.errorf("createFile: unable to close bucket %q, file %q: %v", d.bucketName, fileName, err)
//		return
//	}
//}

// readFile reads the named file in Google Cloud Storage.
func ReadFile(w http.ResponseWriter, r *http.Request) {

	fileName, ok := r.URL.Query()["image"]

	if !ok {
		log.Println("Url Param 'image' is missing")
		return
	}

	ctx := context.Background()
	cli, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	bucket := cli.Bucket(config.BUCKET_NAME)

	filePath := "/blog-images/" + strings.Join(fileName, "") + ".jpg"
	reader, err := bucket.Object(filePath).NewReader(ctx)
	if err != nil {
		log.Println("readFile: unable to open file from bucket %q, file %q: %v", config.BUCKET_NAME, filePath, err)
		return
	}
	defer reader.Close()
	slurp, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Println("readFile: unable to read data from bucket %q, file %q: %v", config.BUCKET_NAME, filePath, err)
		return
	}

	w.Write(slurp)
}
