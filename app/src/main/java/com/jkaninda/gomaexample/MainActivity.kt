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
        Goma.init(this,"http://192.168.194.98/api/")

        getData()
    }


    private fun getData(){
        val parameters:HashMap<String,String> = HashMap()
        parameters["appKey"]="7906b5f78c716fa8385f71368b4d7c7a"
        parameters["magasin_id"]="1"
        val headers:HashMap<String,String> =HashMap()
        headers["accept"]="application/json"
        headers["Authorization"]="Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJhdWQiOiIxIiwianRpIjoiZmIyZTMyZThkMjJjZWJjMGIwOGI0NjgyYWRiZDMwMDZlYzNhYWExNGNmYmJmYmE5YzA5OWE2YzFhOWFiMGJhZWJlODU3YjdhYTQxNTllODciLCJpYXQiOjE2MTk3MjIyODguMzI4OTU5LCJuYmYiOjE2MTk3MjIyODguMzI4OTY0LCJleHAiOjE2NTEyNTgyODguMzE5MzUxLCJzdWIiOiI0Iiwic2NvcGVzIjpbXX0.YVqdizNYZEor4fj3SFnRbNA0UuIvRRrzezUJsC2tkrB0yfwuuvBoOsjssreTR5OrVjltpb7PRs229adFVxrkryHl2HtWpxjQkIcX1R9VUsrgVtwWkA9bGaUXM6wKAUsAl-1Xo-_P_Ge9G3XdHv3w5iCErtkEByjNYWZ2fUPYwlLqZVg-YEa7HItzQ9vOJ2CG2V3FPIh26Mo1pDzqLJ33FjPBVvfe-0xHgAIhcYi37SVivnHcdsEFy6Q8dEdGCHNyDlIIkZHqrg_zselX9LI7HTukzWMgWmaVpiw9ZXPRXhj9NrYO8IrTe1I5EsM4BsXhjoe818D932c5Q4JPtjKGUWd3zS4Javz_2x4xk_3kSUz2RMMMvRgBm8K7Jg-kwPnnvIQYcRCPRrxP7DB2i4aQvOTaZUy3bUV4YTs7yS_zOhqyKhEi_gGUsLFG0r_eCSV_QmWCUZ0-A3H3SMOr0O9izt2tMfAh_S5a7ufIv98nXBzzH4SBW1jQPXfPET1bwkjjmdyNu9z4C_MD1ayXAlvkr8Hy1R0v2lmdDVmTJVEBL0CTUAr4XiYnZYom8lZ0YZ9AI0mdBLsM-yUke3iWPCdxyLKEovFV1xsRKFakwBvsCtHBJjlt5KA8zKSmbWuekPcQZLlq_rG_elxwkCq4K79TBzAUlM5QaM1XSFVv0nfyyh8"


        Goma.del("products/getall",headers,parameters ,object:OnResponseListener{
            override fun onSuccess(response: String?) {

                responsetv.text=response.toString()
            }

            override fun onError(error: String?) {
                showToast(error.toString())
            }

        })

    }

    private fun showToast(message:String){
      return  Toast.makeText(this, message, Toast.LENGTH_SHORT).show()
    }
    
}