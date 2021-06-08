# Android Kotlin Goma http protocol client Library


[![License: Apache-2.0](https://img.shields.io/badge/License-Apache%202.0-yellow.svg)](http://www.apache.org/licenses/LICENSE-2.0)
## Supported methods

GET, POST, PUT, DELETE

## Download and Import

### Android Studio/Gradle

 - Maven:
 
 ```groovy
	<repositories>
		<repository>
		    <id>jitpack.io</id>
		    <url>https://jitpack.io</url>
		</repository>
	</repositories>
	
	
	

	<dependency>
	    <groupId>com.github.jkanTech</groupId>
	    <artifactId>goma</artifactId>
	    <version>1.0.2</version>
	</dependency>


 ```
 
 - JitPack.io, add `jitpack.io` repositiory and dependency to your `build.gradle`:
 
 ```groovy
    repositories {
        maven {
            url "https://jitpack.io"
        }
    }
	
    dependencies {
	        implementation 'jkaninda:goma:1.0.2'
		}
```

   
### Android Studio

 ```groovy
    repositories {
        maven {
            url "https://jitpack.io"
        }
    }
	
    dependencies {
	        implementation 'jkaninda:goma:1.0.2'
		}
 ```
 ### Manifests
 
```xml
    <uses-permission android:name="android.permission.INTERNET"/>

<application

        android:usesCleartextTraffic="true"

```



### Sample Kotlin Usage 
#### GET Method

```Kotlin
 private fun getUsers(){
 // initialize Goma Library + baseURL
 Goma.init(this,"http://192.168.8.101/api/v1/")
 
         
 //GET Method ,path
   Goma.get("users", object : OnResponseListener {
     
           override fun onSuccess(response: String?) {
             
                   users=users.deserializer(response)
 
                 Toast.makeText(requireContext(), response.toString(), Toast.LENGTH_SHORT).show()
             }
            override fun onError(error: String?) {
                 Toast.makeText(requireContext(), error.toString(), Toast.LENGTH_SHORT).show()
             }
 
         })
 
 
 }

```
#### POST Method

```kotlin
private fun addUser(){
// initialize Goma Library + baseURL
Goma.init(this,"http://192.168.8.101/api/v1/")

val parameters:HashMap<String, String> = HashMap()
        parameters.put("appKey","12345#")
        parameters.put("name","Jonas")
        parameters.put("title","Inform")
        parameters.put("age","18")

//POST Method ,path
  Goma.post("add",parameters, object : OnResponseListener {
    
          override fun onSuccess(response: String?) {

                Toast.makeText(requireContext(), "Successfully added", Toast.LENGTH_SHORT).show()
            }
           override fun onError(error: String?) {
                Toast.makeText(requireContext(), error.toString(), Toast.LENGTH_SHORT).show()
            }

        })


}
    

```
#### PUT Method

```kotlin
private fun add(){

// initialize Goma Library 
  Goma.init(this)

 val parameters: HashMap<String, String> = HashMap()
        parameters["Appkey"] = "12345"
        parameters["name"] = "Bob"
        parameters["title"] = "Infor"
        parameters["age"] = "30"


        Goma.put("http://192.168.8.101/api/v1/adduser", parameters,object :OnResponseListener {
            override fun onSuccess(response: String?) {
        Toast.makeText(this, response.toString(), Toast.LENGTH_SHORT).show()

            }

            override fun onError(error: String?) {

        Toast.makeText(this, response.toString(), Toast.LENGTH_SHORT).show()

            }


        })
}

```

#### DELETE Method

```kotlin
private fun delUser(){
// initialize Goma Library
Goma.init(this)

val parameters:HashMap<String, String> = HashMap()
        parameters.put("appKey","12345#")
        parameters.put("id","4")

        
//DELETE Method ,URL
  Goma.del("http://192.168.43.128/api/delUser",parameters, object : OnResponseListener {
    
          override fun onSuccess(response: String?) {

                Toast.makeText(requireContext(), "Successfully deleted", Toast.LENGTH_SHORT).show()
            }
           override fun onError(error: String?) {
                Toast.makeText(requireContext(), error.toString(), Toast.LENGTH_SHORT).show()
            }

        })


}

```


<h2 id="examples">Examples :eyes:</h2>

Download the [Goma Example App]() or look at the [source code](https://github.com/jkanTech/goma/tree/master/CrudExample).


<br/>
 
## Authors

* **Jonas Kaninda**  - [jkanTech](https://github.com/jkantech)


## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
