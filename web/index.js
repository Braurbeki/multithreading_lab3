require('dotenv').config({path:'.env'})
const express = require("express");
const axios = require("axios")
const app = express();

app.use("/static", express.static("public"));
app.use(express.urlencoded({extended: true}));

app.set("view engine", "ejs");

app.get('/', (req, res) => {
    axios.get(`${process.env.API_URL}/getAll`).then(result => {
        res.render('todo.ejs', {todoTasks: result.data});
    })
    
})

app.post('/', (req, res) => {
    axios.post(`${process.env.API_URL}/put`, {
        message: req.body.content
    })
    res.redirect("/");
})

app.get('/remove/:id', (req, res) => {
    axios.delete(`${process.env.API_URL}/delete/${req.params.id}`)
    res.redirect("/");
})

app.route("/edit/:id")
.get((req, res) => {
    const id = req.params.id
    axios.get(`${process.env.API_URL}/getAll`).then(result => {
        res.render('edit.ejs', {todoTasks: result.data, idTask: id});
    })
})
.post((req, res) => {
    const id = req.params.id
    axios.get(`${process.env.API_URL}/get/${id}`).then(result => {
        axios.post(`${process.env.API_URL}/update`, {
            id: result.data.Id,
            message: req.body.content,
            done: result.data.Done
        }).then(result => {
            res.redirect("/")
        })
    })
})

app.post('/update/:id'), (req, res) => {
    console.log(req.params.id)
}

app.listen(process.env.PORT || 5000, () => console.log("Server is UP"))