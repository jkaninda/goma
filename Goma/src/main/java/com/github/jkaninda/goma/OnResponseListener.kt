package com.github.jkaninda.goma

interface OnResponseListener {
    fun onSuccess(response: String?)
    fun onError(error: String?)
}