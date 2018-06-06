package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func getHTTPResponse(url string, t *testing.T) (rr *httptest.ResponseRecorder) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr = httptest.NewRecorder()
	handler := http.HandlerFunc(getThumbnailHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	return rr
}

func TestOkCall(t *testing.T) {
	// Go over all kinds of widths+heights combinations and make sure result is 200
	for w := 10; w <= 1000; w += 150 {
		for h := 10; h <= 1000; h += 150 {

			url := fmt.Sprintf("/thumbnail?url=http://www.ximagic.com/d_im_lenajpeg/lena_comp.jpg&width=%v&height=%v", w, h)
			// Check a "good flow call"
			rr := getHTTPResponse(url, t)
			// Check the status code is what we expect.
			if status := rr.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
		}
	}
}

func TestErrorUrlNonJpeg(t *testing.T) {
	// Check fail we send non-jpeg
	rr := getHTTPResponse("/thumbnail?url=https://image.fnbr.co/outfit/5ab156b3e9847b3170da0324/png.png&width=500&height=900", t)
	expected := http.StatusInternalServerError
	// Check the status code is what we expect.
	if status := rr.Code; status != expected {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expected)
	}
}

func TestErrorUrlBad(t *testing.T) {
	// Check Bad Url address
	rr := getHTTPResponse("/thumbnail?url=bad_value&width=500&height=900", t)
	expected := http.StatusInternalServerError
	// Check the status code is what we expect.
	if status := rr.Code; status != expected {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expected)
	}
}

func TestErrorUrlMissing(t *testing.T) {
	// Check missing url
	rr := getHTTPResponse("/thumbnail?url2=http://www.ximagic.com/d_im_lenajpeg/lena_comp.jpg&width=500&height=900", t)
	expected := http.StatusInternalServerError
	// Check the status code is what we expect.
	if status := rr.Code; status != expected {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expected)
	}
}

func TestErrorWidthBad(t *testing.T) {
	// check bad width values
	badValues := [5]string{"13-", "_14", "56.3", "z", "5T"}
	for _, element := range badValues {
		url := fmt.Sprintf("/thumbnail?url=http://www.ximagic.com/d_im_lenajpeg/lena_comp.jpg&width=%v&height=900", element)
		// element is the element from someSlice for where we are
		rr := getHTTPResponse(url, t)
		expected := http.StatusInternalServerError
		// Check the status code is what we expect.
		if status := rr.Code; status != expected {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, expected)
		}
	}

}

func TestErrorWidthMissing(t *testing.T) {
	// check  width missing
	rr := getHTTPResponse("/thumbnail?url=http://www.ximagic.com/d_im_lenajpeg/lena_comp.jpg&width12=500&height=900", t)
	expected := http.StatusInternalServerError
	// Check the status code is what we expect.
	if status := rr.Code; status != expected {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expected)
	}
}

func TestErrorHeightBad(t *testing.T) {
	// check bad height values
	badValues := [5]string{"13-", "_14", "56.3", "z", "5T"}
	for _, element := range badValues {
		url := fmt.Sprintf("/thumbnail?url=http://www.ximagic.com/d_im_lenajpeg/lena_comp.jpg&width=500&height=%v", element)
		// element is the element from someSlice for where we are
		rr := getHTTPResponse(url, t)
		expected := http.StatusInternalServerError
		// Check the status code is what we expect.
		if status := rr.Code; status != expected {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, expected)
		}
	}

}

func TestErrorHeightMissing(t *testing.T) {
	// check height missing
	rr := getHTTPResponse("/thumbnail?url=http://www.ximagic.com/d_im_lenajpeg/lena_comp.jpg&width=500&height12=900", t)
	expected := http.StatusInternalServerError
	// Check the status code is what we expect.
	if status := rr.Code; status != expected {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expected)
	}
}
