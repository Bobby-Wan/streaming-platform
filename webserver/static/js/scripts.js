async function showStreamConfig(e){
    e.preventDefault();

    $.get(`/configure`,null, function(data){
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