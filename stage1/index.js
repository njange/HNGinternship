const express = require('express');

const app = express();
const PORT = 3000;

//Middleware to enforce json content-type
app.use((req,res,next) => {
  res.setHeader('Content-Type', 'application/json');
  next();
});

//Routes
app.get('/', (req,res) => {
	res.status(200).json({message: "API is running"});

});

app.get('/health', (req,res) => {
        res.status(200).json({message: "healthy"});

});

app.get('/me', (req,res) => {
        res.status(200).json({name:"James Njange",
			      email:"jamesnjange80@gmail.com",
			      github:"https://github.com/njange"});

});

//start the server
app.listen(PORT, () => {
	console.log(`Server is running on port ${PORT}`);
});
