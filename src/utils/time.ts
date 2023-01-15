// credit to alistair.sh

const birth = Date.parse(new Date("1 June 2008 12:00 EST").toString());

const msSinceBirth = Date.parse(Date()) - birth;

export const age = msSinceBirth / 1000 / 60 / 60 / 24 / 365;
