async function showStreamConfig(e){
    e.preventDefault();

    let sel = document.getElementById('categories');
    if(sel.selectedIndex === 0){
        sel.style.background='#ff5f6e';
        return;
    }
    let title=getElementValue("title");
    if(title.value === ""){
        title.style.background='#ff5f6e';
        return;
    }
    let category=getSelectedCategory();
    $.get(`/configure?category=${category}&title=${title}`,
    null, 
    function(data){
        if(data && data.hasOwnProperty('server-address') && data.hasOwnProperty('stream-key')){
          showConfigurationDiv(data['server-address'],data['stream-key']);
        }
        else{
            console.log('invalid response from server')
        }
      },"json").fail(function(response){
        console.log('error');
        console.log(response);
      });
};

function showConfigurationDiv(address, key){
    let addrInput = document.getElementById('server-address');
    let keyInput = document.getElementById('stream-key');
    let div = document.getElementById('configuration');
    let customDiv = document.getElementById('category-picker');

    if(!addrInput || !keyInput || !div || !customDiv){
        debugger;
        console.error('crutial html elements missing on page');
        return;
    }

    addrInput.value = address;
    keyInput.value = key;
    div.classList.remove('invisible');
    customDiv.classList.remove('visible-block');
    customDiv.classList.add('invisible');
};

function getElementValue(id){
    return document.getElementById(id).value;
}

function getSelectedCategory(){
    let e = document.getElementById('categories');
    return e.options[e.selectedIndex].text;
}

function showMessage(message, isAlert){
    span = document.createElement("span");
    span.innerHTML = message;
    span.style.padding = "1.3px";
    if(isAlert){
        span.style.background = 'salmon';
    }
    else{
        span.style.background = 'greenyellow';
    }
    document.getElementsByTagName('body')[0].childNodes.addAt(0, span);
    setTimeout(()=>{
        span.remove();
      }, 3000);
}