Dressformer assets server. File API
============================

### /
_Upload files to asset._

__Endpoint:__ `http://v2.dressformer.com/assets/`  
  
__Methods:__

- POST

__Expected data:__

- MultiPart Form Data

__Example:__   
Following `curl` command uploads selected files 

```sh
curl                                                                                       \
	-i -XPOST -H "ContentType:multipart/form-data"                                         \
	-F name=Base.obj                         -F filedata=@Base.obj                         \
	-F name=Chest_max.obj                    -F filedata=@Chest_max.obj                    \
	-F name=KPL_201407_0020_0005_diffuse.jpg -F filedata=@KPL_201407_0020_0005_diffuse.jpg \
	-F name=KPL_201407_0020_0005_normal.jpg  -F filedata=@KPL_201407_0020_0005_normal.jpg  \
http://v2.dressformer.com/assets/
```
and returns

```json
[
	{
		"id"        : "53f622eb0000000000000001",
		"orig_name" : "Base.obj"
	},
	{
		"id"        : "53f622eb0000000000000002",
		"orig_name" : "Chest_max.obj"
	},
	{
		"id"        : "53f735c10000000000000001",
		"orig_name" : "KPL_201407_0020_0005_diffuse.jpg"
	},
	{
		"id"        : "53f735c10000000000000002",
		"orig_name" : "KPL_201407_0020_0005_normal.jpg"
	}	
]	
```

### /:id or /?geom_ids=
_Get file(s) from asset. Uniform method for every content types._

__Endpoint:__ `http://v2.dressformer.com/assets/`  
  
__Methods:__

- GET

__Expected parameters:__

- Geometry UUID or comma-separated list via `geom_ids` parameter
	- __Parameters:__
		- `height`    (float)
		- `chest`     (float)
		- `underbust` (float)
		- `waist`     (float)
		- `hips`      (float)
	- __Result__:  
		All files will be returned in sequence in same session.  
		Content type `application/octet-stream` and custom header `Df-Sizes` will be set.  
		If header `Accept: application/json` is provided, result will contain array of corresponding geometry objects, content type will be `application/json`.  
	
- ObjectId  
	__Parameters:__

	- `scale` Scaling image to dimensions  
		Prototype: `([0-9]+x) or (x[0-9]+) or ([0-9]+) or (0.[0-9]+)`  
		E.g.:  
  			+ `800x` scale to width 800px, height will be calculated  
		  	+ `x600` scale to height 600px, width will be calculated  
		  	+ `640`  maximum dimension is 640px, e.g. original 1024x768 pixel image will be scaled
  		  	   to 640x480, same option applied for 900x1600 image results 360x640  
		  	+ `0.5`  50% of original dimensions, e.g. 1024x768 = 512x384
	- `q` 0-100 image quality
	- `format` Image format — `png` or `jpg`. Jpeg is default one.
  
__Examples:__

Getting two geometries

```sh
curl -XGET http://v2.dressformer.com/?geom_ids=5ca78b9a-23ed-4551-a6f2-9e3bb9f7c919,b1f8be54-8310-4962-be22-f2446653ea1b

HTTP/1.1 200 OK
Server: nginx/1.4.6 (Ubuntu)
Date: Wed, 29 Oct 2014 15:11:27 GMT
Content-Type: application/octet-stream
Df-Sizes: 1528317,1665655

... some obj data ...

```

Getting some image

```sh
curl -XGET http://v2.dressformer.com/assets/5451082c0000020000000001

HTTP/1.1 200 OK
Server: nginx/1.4.6 (Ubuntu)
Date: Wed, 29 Oct 2014 15:31:24 GMT
Content-Type: image/png

... original png ...

```

Getting some image, scaled to 50% of its dimensions

```sh
curl -XGET http://v2.dressformer.com/assets/5451082c0000020000000001?scale=0.5

HTTP/1.1 200 OK
Server: nginx/1.4.6 (Ubuntu)
Date: Wed, 29 Oct 2014 15:31:24 GMT
Content-Type: image/png

... original png ...

```

Getting some image, converted to jpeg with 70% quality

```sh
curl -XGET -I http://v2.dressformer.com/assets/5451082c0000020000000001?format=jpg&q=70

HTTP/1.1 200 OK
Server: nginx/1.4.6 (Ubuntu)
Date: Wed, 29 Oct 2014 15:31:24 GMT
Content-Type: image/jpg

... jpg data ...

```

### /geometry

List of all geometries.  

__Endpoint:__ `http://v2.dressformer.com/assets/geometry`  

__Methods:__ 

- GET

__Result:__ Array of geometry objects  

__Geometry Object__  

```json
{
	"id" : "e9d7db56-b032-4857-a4dc-7b78e99e86d0",
	"base" : {
		"id" : "544f86eb0000000000000058"
	},
	"name" : "Юбка garment geometry",
	"morph_targets" : [
		{
			"section": "height"
		},
		{
			"section": "chest",
			"sources": [
				{
					"id": "540af7680000000000000002",
					"weight": 105.073
				},
				{
					"id": "540af7680000000000000003",
					"weight": 78.68
				}
			]
		},		
		{
			"section": "underbust"
		},
		{
			"section": "waist"
		},
		{
			"section": "hips"
		}
	]
}

```

### /geometry/:id

Geometry object control.  

__Endpoint:__ `http://v2.dressformer.com/assets/geometry`  

__Methods:__ 

- POST, PUT, DELETE

__Result:__ Geometry objects  





