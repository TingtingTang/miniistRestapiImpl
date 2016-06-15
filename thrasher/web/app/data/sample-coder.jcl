//WORKLOAD  JOB  MSGCLASS=A,SYSUID={{theRunUser}}
//STEP010   EXEC PGM=IEFBR14
// COMMAND 'ROUTE {{theMachine}}, START {{theEnvScript}}'
// COMMAND 'ROUTE {{theMachine}}, START {{theSetupScript}}'
// COMMAND 'ROUTE {{theMachine}}, START {{theExeScript}}'
// COMMAND 'ROUTE {{theMachine}}, START {{theCleanScript}}'