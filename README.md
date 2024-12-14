### 来自[dtm](https://dtm.pub/ref/gozero.html)的灵感，[saga基于redis分布式锁的分布式事务，以go-zero为列](https://github.com/H-H1/SAGA-go-zero)

##### 压测为4000-5000qps下，98%的最终一致性也是可以保证的

> 常见的CAP理论中，C 一致性，A 可用性，P 分区容错性，一般保证P和A，舍弃C，以及最终一致性，本文实现了AP和最终一致性。并给补偿操作有重试机制。
>
> 之前也是学习了DTM，但是发现一个问题，在高并发情况下，dtm事务并不稳定，大量补偿操作不成功，已过多次验证，所以有了本文，本文如下时序图
>
> ![image-20241214202609524](https://github.com/user-attachments/assets/cd42f34d-a3d8-4f31-823e-a12eaddfa56e)

- 这里我准备一个每秒1000并发的dtm压测，

- stock的库存为100个，成功执行的话会减1库存，回滚补偿就会加1

![image-20241214200601723](https://github.com/user-attachments/assets/90856cd0-8150-4a9a-aa9a-1ca9bb7148cd)


- order这里插入成功的话就是新建一条记录，然后row_state的值是 0，回滚补偿了row_state是-1

![image-20241214200354019](https://github.com/user-attachments/assets/742dbe3d-e76c-4ce1-a7d3-96ebde8f2d09)


- 看压测结果，成功217个

![image-20241214200534660](https://github.com/user-attachments/assets/b6999b00-6ac4-4037-9f3a-2efe48a30e72)


- stock没有问题，成功扣除完成

![image-20241214200706228](https://github.com/user-attachments/assets/46e99d56-f299-4606-8096-0016a7c74ae5)


- 但是order出问题了，插入了大量的成功状态的记录，大量补偿失败了

![image-20241214201343030](https://github.com/user-attachments/assets/4a8dd73a-72c7-4416-b17b-78c308278b65)


#### 于是我打算自己写一个基于redis分布式锁的分布式事务，防止补偿失败的问题。

这个事务目前可以提供2个微服务的一致性，需要2个以上的我感觉可以先合并下服务

具体时序图下

![image-20241214202609524](https://github.com/user-attachments/assets/cd42f34d-a3d8-4f31-823e-a12eaddfa56e)


采用两个err分别表示确认微服务的状态

1. Create(l.ctx, createOrderReq)如果失败了，服务直接返回失败，啥没发生
2. Create(l.ctx, createOrderReq)成功了，row_state为0，但是l.svcCtx.StockRpc.Deduct(l.ctx, deductReq)失败，就会进行补偿order，补偿加入了重试
3. 两个微服务都加入了分布式锁+uid的方法，防止延迟导致的误删。

- 进行压测，可以看到由于锁的存在1，成功数只有4了

![image-20241214204252365](https://github.com/user-attachments/assets/f2de42e3-ce75-468f-ace1-5d73c80c8b38)


但是，order和stock不在出现补偿失败导致的不一致的情况。扣除4个，4个0状态

![image-20241214204419091](https://github.com/user-attachments/assets/bb53d5ea-09ce-4a54-aca5-0f694326c450)

![image-20241214204433122](https://github.com/user-attachments/assets/1909edd3-0ee3-4282-b652-391618e014e1)



多次压测的结果一致。

![image-20241214204604732](https://github.com/user-attachments/assets/975df20c-76b4-42da-b80a-0eb13b881130)

![image-20241214204613884](https://github.com/user-attachments/assets/e748435e-fc9a-488f-95c0-176518141234)


![image-20241214204626565](https://github.com/user-attachments/assets/2046eedb-c117-41f5-8720-07246d1c0946)



压测为4000-5000qps下，98%的最终一致性也是可以保证的

![image-20241214210804031](https://github.com/user-attachments/assets/f6485935-cff3-481c-b66a-7b2ba6aaebb1)


![image-20241214210811660](https://github.com/user-attachments/assets/ab54f6f0-331b-4fa9-b441-a2131b722b6d)

![image-20241214210819236](https://github.com/user-attachments/assets/78009ccb-2af8-4793-9243-9bf21cf0b445)

