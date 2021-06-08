package com.github.jkaninda.goma.Serializer

import com.google.gson.GsonBuilder
import org.json.JSONObject
import java.io.BufferedReader
import java.io.FileReader
import java.io.FileWriter
import java.io.IOException

@Suppress("NULLABILITY_MISMATCH_BASED_ON_JAVA_ANNOTATIONS")
class Serializer {
    private val gson =
            GsonBuilder().enableComplexMapKeySerialization().serializeNulls().setPrettyPrinting()
                    .create()

    fun Serializer() {}

    fun <T> serialize(filename: String?, objectTString: String?, classOfT: Class<T>?): Boolean {
        val p = gson.fromJson(objectTString, classOfT)
        return this.serialize(filename, p, classOfT)
    }

    fun <T> serialize(filename: String?, objectOfT: T, classOfT: Class<T>?): Boolean {
        return try {
            val w = FileWriter(filename)
            w.write(gson.toJson(objectOfT))
            w.close()
            true
        } catch (var6: IOException) {
            false
        }
    }

    fun <T> deserializeToString(filename: String?, classOfT: Class<T>?): String? {
        val obj: T? = deserialize(filename, classOfT)
        return gson.toJson(obj)
    }

    fun <T> deserialize(filename: String?, classOfT: Class<T>?): T? {
        return try {
            val reader = BufferedReader(FileReader(filename))
            gson.fromJson(reader, classOfT)
        } catch (var4: Exception) {
            null
        }
    }

    fun <T> toJson(objectOfT: T): String? {
        return gson.toJson(objectOfT)
    }


    fun <T> fromJson(objectTString: String?, classOfT: Class<T>?): T {
        return gson.fromJson(objectTString, classOfT)
    }



    fun fetchJson( jsonString:String?): JSONObject {
        if (jsonString!=null){

         return   JSONObject(jsonString)
            
        }
        return JSONObject("null")
        }


}