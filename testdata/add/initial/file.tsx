import React from 'react';

interface HelloProps {
    name: string;
}

const HelloWorld: React.FC<HelloProps> = ({ name }) => {
    return (
        <div>
            <h1>Hello, {name}!</h1>
        </div>
    );
};

export default HelloWorld;