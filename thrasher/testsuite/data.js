// mgo insert nested arrays, code example:
//     http://stackoverflow.com/questions/32652182/mgo-how-to-find-nested-document-inside-nested-arrays
// 
{
  _id,
  name:    "uss_kernel",
  creator: "mike",
  modifier: "mike",
  create_time: date,
  update_time: date,
  team:    "uss",
  setup_script: "submit setup.jcl",
  clean_script: "submit back_out.jcl",
  command: "S kernel_case",
  testcases: [
                 { 
                   testcase: "name",
                   caseid: _id,
                 }
             ]
  machines: [ np0 np1 np2 np3 ],
  elapsed_time: 3600,
  desc: "OS kernel case"
}
