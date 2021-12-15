// 리액트 애플리케이션의 메인 컴포넌트로, 다른 모든 컴포넌트를 결합합니다.
// 이 파일의 App 컴포넌트는 사용자 로그인과 로그아웃 등의 중요한 작업을 처리합니다.
import React from 'react';
import CardContainer from './ProductCards';
import { BrowserRouter as Router, Route } from "react-router-dom";
import Nav from './Navigation';
import { SignInModalWindow, BuyModalWindow } from './modalwindows';
import About from './About';
import Orders from './orders';
import Cookies from 'js-cookie';


class App extends React.Component {

  constructor(props) {
    super(props);
    const user = Cookies.get("user") || {loggedin:false};
    this.state = {
      user: user,
      showSignInModal: false,
      showBuyModal: false
    };
    this.handleSignedIn = this.handleSignedIn.bind(this);
    this.handleSignedOut = this.handleSignedOut.bind(this);
    this.showSignInModalWindow = this.showSignInModalWindow.bind(this);
    this.toggleSignInModalWindow = this.toggleSignInModalWindow.bind(this);
    this.showBuyModalWindow = this.showBuyModalWindow.bind(this);
    this.toggleBuyModalWindow = this.toggleBuyModalWindow.bind(this);
    
  }

  handleSignedIn(user) {
    console.log("Sign in happening...");
    const state = this.state;
    const newState = Object.assign({},state,{user:user,showSignInModal:false});
    this.setState(newState);
  }

  handleSignedOut(){
    console.log("Call app signed out...");
    const state = this.state;
    const newState = Object.assign({},state,{user:{loggedin:false}});
    this.setState(newState);
    Cookies.set("user",{loggedin:false});
  }

  showSignInModalWindow(){
    const state = this.state;
    const newState = Object.assign({},state,{showSignInModal:true});
    this.setState(newState);
  }

  toggleSignInModalWindow() {
    const state = this.state;
    const newState = Object.assign({},state,{showSignInModal:!state.showSignInModal});
    this.setState(newState);
  }



  showBuyModalWindow(id,price){
    const state = this.state;
    const newState = Object.assign({},state,{showBuyModal:true,productid:id,price:price});
    this.setState(newState);
  }

  toggleBuyModalWindow(){
    const state = this.state;
    const newState = Object.assign({},state,{showBuyModal:!state.showBuyModal});
    this.setState(newState); 
  }
  //location='user.json'
  render() {
    return (
      <div>
        <Router>
          <div>
            <Nav user={this.state.user} handleSignedOut={this.handleSignedOut} showModalWindow={this.showSignInModalWindow}/>
            <div className='container pt-4 mt-4'>
              <Route exact path="/" render={() => <CardContainer location='/products' showBuyModal={this.showBuyModalWindow}/>} />
              <Route path="/promos" render={() => <CardContainer location='/promos' promo={true} showBuyModal={this.showBuyModalWindow}/>} />
              {this.state.user.loggedin ? <Route path="/myorders" render={()=><Orders location={'/user/'+this.state.user.ID+'/orders'}/>}/> : null}
              <Route path="/about" component={About} />
            </div>
            <SignInModalWindow handleSignedIn={this.handleSignedIn} showModal={this.state.showSignInModal} toggle={this.toggleSignInModalWindow} />
            <BuyModalWindow showModal={this.state.showBuyModal} toggle={this.toggleBuyModalWindow} user={this.state.user.ID} productid={this.state.productid} price={this.state.price}/>
          </div>
        </Router>
      </div>
    );
  }
}
export default App;