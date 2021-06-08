package com.github.jkaninda.goma.internal


import com.android.volley.Response
import com.android.volley.toolbox.StringRequest


interface RequestFactory {
    class HttpRequest(
            method: Int, URL: String?,private val headers:MutableMap<String,String>,private val params: MutableMap<String, String>,
            listener: Response.Listener<String?>?,
            error: Response.ErrorListener?
    ) : StringRequest(method, URL, listener, error) {


        override fun getParams(): MutableMap<String, String> {
            return params
        }

        override fun getHeaders(): MutableMap<String, String> {
            return headers
        }


    }
}




