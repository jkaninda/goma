package com.github.jkaninda.goma

import android.content.Context
import android.util.Log
import com.android.volley.Request
import com.android.volley.toolbox.StringRequest
import com.android.volley.toolbox.Volley
import com.github.jkaninda.goma.internal.RequestFactory

object Goma {
    private lateinit var appContext: Context
    private var URL: String = ""
    private val queue by lazy { Volley.newRequestQueue(appContext) }
    private var params: HashMap<String, String> = HashMap()
    private var headers: HashMap<String, String> = HashMap()
    private var TAG="Goma"


    /**
     * @param context The context to use.
     *                 or android.app.Activity} object.
     * @param baseURL Server URL
     */
    fun init(context: Context,baseURL:String?=null) {
        if (!::appContext.isInitialized) {
            appContext = context.applicationContext
            if (baseURL != null) {
                if (baseURL.isNotEmpty()) {
                    URL = if (baseURL.last().toString() == "/") {
                        baseURL
                    } else {
                        "$baseURL/"
                    }
                }
            }

        }
    }
/**
 * @param path path
 * @param headers HTTP headers
 * @param parameters parameters
 * @param listener Object:OnResponseListener({
 * })
 */
    fun get(path: String, parameters: HashMap<String, String>? = HashMap(), listener: OnResponseListener?) {

        val req = StringRequest(
            Request.Method.GET,
            "$URL${checkPath(path)}",
            { p1 -> listener?.onSuccess(p1) }) { p1 ->
            if (listener != null) {
                listener.onError(p1.message)
                p1.printStackTrace()
            }

        }

        queue.add(req)
    }





    /**
     * @param path path
     * @param headers HTTP headers
     * @param parameters parameters
     * @param listener Object:OnResponseListener({
     * })
     */
    fun post(
        path: String, headers:HashMap<String,String>,
        parameters: HashMap<String, String>? = null,
        listener: OnResponseListener?
    ) {
        if (parameters != null) {
            params = parameters
        }
        postRequest("$URL${checkPath(path)}", headers, params, listener)

    }


    /**
     * @param path path
     * @param headers HTTP headers
     * @param parameters parameters
     * @param listener Object:OnResponseListener({
     * })
     */
    fun del(
        path: String, headers: HashMap<String, String>,
        parameters: HashMap<String, String>? = null,
        listener: OnResponseListener?
    ) {
        if (parameters != null) {
            params = parameters
        }
        delRequest("$URL${checkPath(path)}",headers, params, listener)

    }

    /**
     * @param path path
     * @param headers HTTP headers
     * @param parameters parameters
     * @param listener Object:OnResponseListener({
     * })
     */
    fun put(
        path: String,headers:HashMap<String,String>,
        parameters: HashMap<String, String>? = HashMap(),
        listener: OnResponseListener?
    ) {
        if (parameters != null) {
            params = parameters
        }
        putRequest("$URL${checkPath(path)}",headers, params, listener)


    }/**
     * @param path path
     * @param headers HTTP headers
     * @param parameters parameters
     * @param listener Object:OnResponseListener({
     * })
     */
    fun patch(
        path: String,headers:HashMap<String,String>,
        parameters: HashMap<String, String>? = HashMap(),
        listener: OnResponseListener?
    ) {
        if (parameters != null) {
            params = parameters
        }
        patchRequest("$URL${checkPath(path)}",headers, params, listener)


    }


    /*--------------------------------------------CANCEL----------------------------------------------*/

    fun cancel() {
        queue.cancelAll(TAG)
    }

    fun stop() {
        Log.d(TAG,"Stoped")
        queue.stop()
    }


    private fun postRequest(
        path: String,headers:HashMap<String,String>,
        parameters: HashMap<String, String>,
        listener: OnResponseListener?
    ) {
        val req = RequestFactory.HttpRequest(
            Request.Method.POST,
            path,headers,
            parameters,
            { p1 -> listener?.onSuccess(p1) }) { p1 ->
            if (listener != null) {
                listener.onError(p1.message)
                p1.printStackTrace()
            }
        }
        queue.add(req)
    }


    private fun delRequest(
        path: String,headers:HashMap<String,String>,
        parameters: HashMap<String, String>,
        listener: OnResponseListener?
    ) {
        val req = RequestFactory.HttpRequest(
            Request.Method.DELETE,
            path,headers,
            parameters,
            { p1 -> listener?.onSuccess(p1) }) { p1 ->
            if (listener != null) {
                listener.onError(p1.message)
                p1.printStackTrace()
            }
        }
        queue.add(req)
    }

    private fun putRequest(
        path: String,headers:HashMap<String,String>,
        parameters: HashMap<String, String>,
        listener: OnResponseListener?
    ) {
        val req = RequestFactory.HttpRequest(Request.Method.PUT, path,headers, parameters, { p1 -> listener?.onSuccess(p1) }) { p1 ->
            if (listener != null) {
                listener.onError(p1.message)
                p1.printStackTrace()
            }
        }
        queue.add(req)
    }

    private fun patchRequest(
        path: String, headers:HashMap<String,String>,
        parameters: HashMap<String, String>,
        listener: OnResponseListener?
    ) {
        val req = RequestFactory.HttpRequest(Request.Method.PATCH, path,headers, parameters, { p1 -> listener?.onSuccess(p1) }) { p1 ->
            if (listener != null) {
                listener.onError(p1.message)
                p1.printStackTrace()
            }
        }
        queue.add(req)
    }

    private fun checkPath(path: String?): String {
        if (path!!.first().toString() == "/") {


            return path.removePrefix("/")
        }
        return path

    }



}



















