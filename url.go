package main

import(
"fmt"
"net/url"
"net/http"
"log"
"gopkg.in/mgo.v2"
"gopkg.in/mgo.v2/bson"
"github.com/gorilla/mux"
"github.com/speps/go-hashids"
//"go.mongodb.org/mongo-driver/mongo"
)

const(
  database = "go"
  collection = "urls"
  col = "user"
)

type Url struct{
Id  int  `json:"_id,omitempty" bson:"_id,omitempty"`
ShortUrl string
LongUrl  string
Count int
}

func init(){
 session, err := mgo.Dial("localhost")
 if err != nil {
                panic(err)
        }
defer session.Close()
}


func main(){
        r := mux.NewRouter()
        r.HandleFunc("/",welcome)
        r.HandleFunc("/UrlValid",UrlValid)
        if err := http.ListenAndServe(":80", r); err != nil {
                log.Fatal(err)
         }
}

func IsUrl(str string) bool {
    u, err := url.Parse(str)
    return err == nil && u.Scheme != "" && u.Host != ""
}

func welcome(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w,htmlStr)
}

var ul Url
var ip = "104.154.237.49"
func UrlValid(w http.ResponseWriter, r *http.Request){
     r.ParseForm()
    session, err := mgo.Dial("localhost")
       if err != nil {
                panic(err)
        }

     c := session.DB(database).C(collection)
     ul.LongUrl=r.FormValue("name")
     fmt.Println(ul.LongUrl)
     if (IsUrl(ul.LongUrl))==true{
//        fmt.Fprintln(w,"Creating ShortUrl")
          k:= check(ul)
          if k==false{
          c.Insert(&ul)
          fmt.Println("Creating Url")
          q := session.DB(database).C(collection)
     err:= q.Find(bson.M{"longurl":ul.LongUrl}).One(&ul)
        if err!=nil{
                fmt.Println("error")
             }
             k:=Create([]int {ul.Id,0})
         ul.ShortUrl = k
          q.Update(bson.M{"longurl":ul.LongUrl},&ul)
          } else{
                fmt.Fprintln(w,ul.ShortUrl)
                fmt.Fprintln(w,"the above url is short url for this link")
            } 
      } else{
       fmt.Fprintf(w,"The URL entered is not valid") 
         }
  }


func check(s Url) bool {
   session, err := mgo.Dial("localhost")
       if err != nil {
                panic(err)
        }

     c := session.DB(database).C(collection)
     pipe := c.Pipe([]bson.M{{"$match": bson.M{"longurl":s.LongUrl}}})
     resp :=[]bson.M{}
     k := pipe.All(&resp)
     if k != nil {
     fmt.Println("qqq")
     }
    fmt.Println(resp)
      fmt.Println(len(resp))
   if len(resp)==0{
   //  fmt.Println("Not Present")
    return false
   } else{
    fmt.Println(resp)
    return  true
      }

}

func Create(Obid [] int)( su string){

       hd := hashids.NewData()
        h, _ := hashids.NewWithData(hd)
        e, _ := h.Encode(Obid)
       Short :="http://"+ip+":80/"+e
      return Short
}

var htmlStr = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8" />
</head>
<body>
  <div>
      <form  action="/UrlValid">
      LongUrl: <input type="text" name="name" value="http://google.com" >
          <input type="submit" value="submit" />
      </form>


  </div>
</body>
</html>
`

