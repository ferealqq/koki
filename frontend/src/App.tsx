import {useState, useEffect} from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import {Greet, GetMostPressedKey, GetKeyEventData} from "../wailsjs/go/main/App";
import { Box, Button, Flex, Select, Spacer } from '@chakra-ui/react';

type EventData = {
    Char: string;
    Count: number;
    Value: number;
    ValueName: string;
}

function App() {
    const [resultText, setResultText] = useState("Please enter your name below 👇");
    const [data, setData] = useState<EventData[]>([]);
    const [count, setCount] = useState(0);
    const [char, setChar] = useState('');
    const [name, setName] = useState('');
    const updateName = (e: any) => setName(e.target.value);
    const updateResultText = (result: string) => setResultText(result);

    function greet() {
        Greet(name).then(updateResultText);
    }

    const onKeyPress = async (evt: any) => {
        console.log(evt.target.value)
        await GetKeyEventData(evt.key).then((res) => {
            setData(res)
        })
        // setChar(`code: ${evt.code}, key: ${evt.key}`)
    }

    useEffect(() => {
        window.addEventListener("keypress", onKeyPress);

        (async () => {
            await GetMostPressedKey().then((res:any) =>{
                console.log(res)
                setCount(res.Count)
                setChar(res.Char)
            })
        })()

        return () => window.removeEventListener("keypress", onKeyPress);
    }, [])

    return (
        <Box w="100%">
            <Flex>
                <Box maxW="sm">
                    <Select placeholder='Today'>
                        <option value='1'>Today</option>
                        <option value='7'>Last 7 days</option>
                        <option value='30'>Last 30 days</option>
                    </Select>
                </Box>

                <Spacer/>
                
                <Box maxW="sm">
                    <Button color="white" bg="green">
                        Logger is active                    
                    </Button>
                </Box>
            </Flex>
            {count !== 0 ? "Char "+char+" count: "+count : "Count ei ole vielä saapunut"}
            {data.map(({Value,ValueName,Count,Char}) => {
                return <h1>Value Name {ValueName} Count {Count} Char {Char}</h1>
            })}
            {/* <div id="input" className="input-box">
                <input id="name" className="input" onChange={updateName} autoComplete="off" name="input" type="text"/>
                <button className="btn" onClick={greet}>Greet</button>
            </div> */}
        </Box>

    )
}

export default App
