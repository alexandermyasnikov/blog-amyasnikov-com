//----- (00000001000063C0) ----------------------------------------------------
void __cdecl -[CreateAccounViewController createAccountClicked:](CreateAccounViewController *self, SEL a2, id a3)
{
  UITextField *v4; // x20
  UITextField *v5; // x21
  UITextField *v6; // x23
  NSString *v7; // x21
  NSUserDefaults *v8; // x25
  const char *v9; // x24
  NSString *v10; // x26
  NSString *v11; // x23
  NSString *v12; // x27
  NSString *v13; // x26
  NSURL *v14; // x24
  NSMutableURLRequest *v15; // x25
  NSString *v16; // x26
  NSString *v17; // x26
  NSString *v18; // x27
  id v19; // x26
  NSString *v20; // x19
  NSString *v21; // x27
  NSString *v22; // x27
  NSString *v23; // x19
  NSString *v24; // x27
  NSString *v25; // x19
  NSData *v26; // x27
  id v27; // x28
  NSURLSession *v28; // x19
  NSURLSessionDataTask *v29; // x22
  NSString *v30; // [xsp+10h] [xbp-90h]
  NSString *v31; // [xsp+18h] [xbp-88h]
  __int64 v32[5]; // [xsp+20h] [xbp-80h] BYREF
  id v33; // [xsp+48h] [xbp-58h] BYREF

  v4 = objc_retainAutoreleasedReturnValue(-[CreateAccounViewController usernameField](self, "usernameField", a3));
  v31 = objc_retainAutoreleasedReturnValue(-[UITextField text](v4, "text"));
  objc_release(v4);
  v5 = objc_retainAutoreleasedReturnValue(-[CreateAccounViewController passwordField](self, "passwordField"));
  v30 = objc_retainAutoreleasedReturnValue(-[UITextField text](v5, "text"));
  objc_release(v5);
  v6 = objc_retainAutoreleasedReturnValue(-[CreateAccounViewController refcodeField](self, "refcodeField"));
  v7 = objc_retainAutoreleasedReturnValue(-[UITextField text](v6, "text"));
  objc_release(v6);
  v8 = objc_retainAutoreleasedReturnValue(+[NSUserDefaults standardUserDefaults](&OBJC_CLASS___NSUserDefaults, "standardUserDefaults"));
  v9 = "https://t-galera-u12p5koe.spbctf.net";
  v10 = objc_retainAutoreleasedReturnValue(
          +[NSString stringWithCString:](
            &OBJC_CLASS___NSString,
            "stringWithCString:",
            "https://t-galera-u12p5koe.spbctf.net"));
  v11 = objc_retainAutoreleasedReturnValue(-[NSUserDefaults stringForKey:](v8, "stringForKey:", v10));
  objc_release(v10);
  objc_release(v8);
  if ( v11 )
    v9 = -[NSString cString](objc_retainAutorelease(v11), "cString");
  v12 = objc_retainAutoreleasedReturnValue(+[NSString stringWithCString:](&OBJC_CLASS___NSString, "stringWithCString:", "%s/%s"));
  v13 = objc_retainAutoreleasedReturnValue(+[NSString stringWithFormat:](&OBJC_CLASS___NSString, "stringWithFormat:", v12, v9, "create"));
  v14 = objc_retainAutoreleasedReturnValue(+[NSURL URLWithString:](&OBJC_CLASS___NSURL, "URLWithString:", v13));
  objc_release(v13);
  objc_release(v12);
  v15 = objc_retainAutoreleasedReturnValue(+[NSMutableURLRequest requestWithURL:](&OBJC_CLASS___NSMutableURLRequest, "requestWithURL:", v14));
  v16 = objc_retainAutoreleasedReturnValue(+[NSString stringWithCString:](&OBJC_CLASS___NSString, "stringWithCString:", "POST"));
  -[NSMutableURLRequest setHTTPMethod:](v15, "setHTTPMethod:", v16);
  objc_release(v16);
  v17 = objc_retainAutoreleasedReturnValue(+[NSString stringWithCString:](&OBJC_CLASS___NSString, "stringWithCString:", "application/json"));
  v18 = objc_retainAutoreleasedReturnValue(+[NSString stringWithCString:](&OBJC_CLASS___NSString, "stringWithCString:", "Content-Type"));
  -[NSMutableURLRequest setValue:forHTTPHeaderField:](v15, "setValue:forHTTPHeaderField:", v17, v18);
  objc_release(v18);
  objc_release(v17);
  v19 = objc_alloc_init((Class)&OBJC_CLASS___NSMutableDictionary);
  v20 = objc_retainAutoreleasedReturnValue(+[NSString stringWithCString:](&OBJC_CLASS___NSString, "stringWithCString:", "%s"));
  v21 = objc_retainAutoreleasedReturnValue(+[NSString stringWithFormat:](&OBJC_CLASS___NSString, "stringWithFormat:", v20, "username"));
  objc_msgSend(v19, "setObject:forKey:", v31, v21);
  objc_release(v21);
  objc_release(v20);
  v22 = objc_retainAutoreleasedReturnValue(+[NSString stringWithCString:](&OBJC_CLASS___NSString, "stringWithCString:", "%s"));
  v23 = objc_retainAutoreleasedReturnValue(+[NSString stringWithFormat:](&OBJC_CLASS___NSString, "stringWithFormat:", v22, "password"));
  objc_msgSend(v19, "setObject:forKey:", v30, v23);
  objc_release(v23);
  objc_release(v22);
  if ( -[NSString length](v7, "length") )
  {
    v24 = objc_retainAutoreleasedReturnValue(+[NSString stringWithCString:](&OBJC_CLASS___NSString, "stringWithCString:", "%s"));
    v25 = objc_retainAutoreleasedReturnValue(+[NSString stringWithFormat:](&OBJC_CLASS___NSString, "stringWithFormat:", v24, "ref_code"));
    objc_msgSend(v19, "setObject:forKey:", v7, v25);
    objc_release(v25);
    objc_release(v24);
  }
  v33 = 0LL;
  v26 = objc_retainAutoreleasedReturnValue(
          +[NSJSONSerialization dataWithJSONObject:options:error:](
            &OBJC_CLASS___NSJSONSerialization,
            "dataWithJSONObject:options:error:",
            v19,
            0LL,
            &v33));
  v27 = objc_retain(v33);
  if ( v26 )
  {
    -[NSMutableURLRequest setHTTPBody:](v15, "setHTTPBody:", v26);
    v28 = objc_retainAutoreleasedReturnValue(+[NSURLSession sharedSession](&OBJC_CLASS___NSURLSession, "sharedSession"));
    v32[0] = (__int64)_NSConcreteStackBlock;
    v32[1] = 3254779904LL;
    v32[2] = (__int64)sub_10000685C;
    v32[3] = (__int64)&unk_100010158;
    v32[4] = (__int64)self;
    v29 = objc_retainAutoreleasedReturnValue(
            -[NSURLSession dataTaskWithRequest:completionHandler:](
              v28,
              "dataTaskWithRequest:completionHandler:",
              v15,
              v32));
    -[NSURLSessionDataTask resume](v29, "resume");
    objc_release(v29);
    objc_release(v28);
  }
  objc_release(v26);
  objc_release(v27);
  objc_release(v19);
  objc_release(v15);
  objc_release(v14);
  objc_release(v11);
  objc_release(v7);
  objc_release(v30);
  objc_release(v31);
}