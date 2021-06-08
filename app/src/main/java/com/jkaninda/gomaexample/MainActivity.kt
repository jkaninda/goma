package com.jkaninda.gomaexample

import androidx.appcompat.app.AppCompatActivity
import android.os.Bundle
import android.widget.TextView
import android.widget.Toast
import com.github.jkaninda.goma.Goma
import com.github.jkaninda.goma.OnResponseListener

class MainActivity : AppCompatActivity() {
    lateinit var responsetv:TextView
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)
        responsetv=findViewById(R.id.response)

    }


}