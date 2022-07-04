import {useState, useEffect} from 'react';
import './App.css';
import {GetMostPressedKey, GetKeyEventData, IsLoggerActive, ToggleLoggerDaemon, GetKeysPressedIn} from "../wailsjs/go/main/App";
import { Box, Text, Flex, Icon, IconButton, Select, Spacer, useToast } from '@chakra-ui/react';
import { MdStop, MdNotStarted } from 'react-icons/md';

type EventData = {
    Char: string;
    Count: number;
    Value: number;
    ValueName: string;
}

function App() {
    const toast = useToast();
    const [data, setData] = useState<EventData[]>([]);
    const [count, setCount] = useState(0);
    const [char, setChar] = useState('');

    const [isActive, setActive] = useState(false);

    const onKeyPress = async (evt: any) => {
        await GetKeyEventData(evt.key).then((res) => {
            setData(res)
        })
    }

    useEffect(() => {
        window.addEventListener("keypress", onKeyPress);

        (() => {
            GetMostPressedKey().then((res:any) => {
                setCount(res.Count)
                setChar(res.Char)
            })
        })();

        (() => {
            IsLoggerActive().then(res => {
                if(typeof res === "boolean") 
                    setActive(res)
                else {
                    // FIXME: set error to true 
                    setActive(false)
                } 
            })
        })();

        return () => window.removeEventListener("keypress", onKeyPress);
    }, [])

    const toggleLogger = () => {
        ToggleLoggerDaemon().then(res => {
            if(typeof res === "boolean"){
                IsLoggerActive().then(active => {
                    let title, description;
                    if(active){
                        title = "Logger deamon successfully started!";                        
                        description = "Something";
                    }else{
                        title = "Logger daemon successfully stopped!";
                        description = "Something";
                    }
                    toast({
                        title,
                        description,
                        status: 'success',
                        duration: 9000,
                        isClosable: true,
                    });
                    setActive(active)
                })
            }else{
                toast({
                    title: 'Something went wrong while toggling logger daemon',
                    description: "Navigate to the logs of koki to find the reason of deamon action failure.",
                    status: 'warning',
                    duration: 9000,
                    isClosable: true,
                })
                throw res;
            }
        });
    }

    return (
        <Box w="100%">
            <Flex>
                <Box maxW="sm" p="2">
                    <Select placeholder='Today'>
                        <option value='1'>Today</option>
                        <option value='7'>Last 7 days</option>
                        <option value='30'>Last 30 days</option>
                    </Select>
                </Box>

                <Spacer/>

                <Box maxW="sm" display="flex" justifyContent="center" alignContent="center" p="2">
                    {
                        isActive ? 
                            <IconButton
                                aria-label='Stop the logger daemon'
                                onClick={toggleLogger}
                                bg="red.400"
                                m="1"
                                icon={
                                    <Icon as={MdStop} boxSize={10} color="white"></Icon>
                                }></IconButton>
                        :
                            <IconButton
                                aria-label='Stop the logger daemon'
                                onClick={toggleLogger}
                                bg="green.400"
                                m="1"
                                icon={
                                    <Icon as={MdNotStarted} boxSize={10} color="white"></Icon>
                                }></IconButton>
                    }
                    <Spacer></Spacer>
                    {
                        isActive ? 
                                <Text color="grey" my="auto" mx="1">
                                    Logger is active                    
                                </Text>
                            :
                                <Text mx="1" my="auto" color="grey">
                                    Logger is not active
                                </Text>                         
                    }
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
            <Today />
        </Box>
    )
}

type HourEvents = {
    Hour: number;
    Count: number;
}

const Today = () => {
    const [hourEventCount, setHEC] = useState<HourEvents[]>([]);
    
    useEffect(()=>{
        (() => {
            GetKeysPressedIn(12).then(data => {
                setHEC(data)
            })
        })()
    },[])

    // return ()
    return <Text> size: {hourEventCount.length} </Text>
}

export default App
